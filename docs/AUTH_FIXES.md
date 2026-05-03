# Authentication and CORS Fixes

## Overview
This document describes the fixes implemented to resolve authentication and CORS issues in the child-bot application.

## Issues Fixed

### 1. Child Profile ID Null Issue
**Problem**: `childProfileId` was always null after onboarding because the authentication flow didn't properly persist and reuse the profile ID.

**Root Cause**: 
- After profile creation during onboarding, the `childProfileId` was stored in `vkStorage` but not in the global auth cache (`sessionStorage` used by `auth.ts`)
- The `getCurrentChildProfileId()` function queried the backend but didn't reuse the cached ID from `vkStorage`
- API requests didn't automatically include `X-Child-Profile-ID` header

**Solution**:
- Updated `frontend/src/lib/auth.ts` to:
  - Check `sessionStorage` for cached profile ID first (fast path)
  - Properly cache profile ID after query or creation
  - Added `createAuthenticatedClient()` function that automatically adds authentication headers
  - Added `getCurrentProfileIdOrFallback()` for synchronous access

- Updated `frontend/src/pages/Onboarding/OnboardingPageNew.tsx` to:
  - Import and call `setCurrentChildProfileId()` after profile creation
  - Ensure immediate availability of profile ID for subsequent requests

### 2. CORS Mixed Content Issue
**Problem**: CORS errors when mixing HTTP and HTTPS requests between frontend and API.

**Root Cause**:
- Strict origin checking rejected requests from different protocols (HTTP vs HTTPS)
- Development mode didn't properly handle mixed content scenarios
- Missing `X-VK-Sign` header support for VK Mini Apps signature validation

**Solution**:
- Updated `api/internal/api/middleware/cors.go` to:
  - Allow both HTTP and HTTPS origins for development/staging
  - Add `Access-Control-Expose-Headers` for client access to auth headers
  - Add `X-VK-Sign` to allowed headers
  - Log origin allowance in development mode

- Updated `.env.production`:
  - Added both `http://` and `https://` versions of allowed origins
  - Added localhost origins for development

### 3. DOM Validation Errors (Nested Buttons)
**Problem**: Original report mentioned nested button elements causing DOM validation errors.

**Finding**: Upon inspection of `OnboardingPageNew.tsx`, no nested `<button>` elements were found. The code uses:
- `<button>` elements containing `<div>` checkboxes (not nested buttons)
- Proper structure with `stopPropagation()` on inner clickable elements

**Status**: No DOM validation issues found; code structure is correct.

## Files Modified

### Frontend
1. **`frontend/src/lib/auth.ts`** - Enhanced authentication utilities
   - Added `createAuthenticatedClient()` for automatic header injection
   - Improved caching strategy with sessionStorage
   - Better error logging
   - Added fallback functions

2. **`frontend/src/pages/Onboarding/OnboardingPageNew.tsx`** - Fixed profile persistence
   - Added `setCurrentChildProfileId()` call after profile creation
   - Ensures immediate availability for subsequent requests

### Backend
3. **`api/internal/api/middleware/cors.go`** - Enhanced CORS handling
   - Support for mixed HTTP/HTTPS origins
   - Added exposed headers
   - Added `X-VK-Sign` header support
   - Development mode logging

4. **`.env.production`** - Updated allowed origins
   - Added both HTTP and HTTPS variants
   - Added localhost for development

## API Authentication Flow

### Protected Endpoints
All endpoints except the following require authentication:
- `/health` - Health check
- `/onboarding/*` - Onboarding endpoints
- `/avatars` - Avatar data
- `/analytics/events` - Analytics
- `/legal/*` - Legal documents
- `/webhooks/*` - Webhook endpoints

### Authentication Headers
For protected endpoints, include:
```
X-Platform-ID: vk
X-Child-Profile-ID: <uuid>
```

For VK-signed requests:
```
X-VK-Sign: <signature>
```

### Example Protected Request
```javascript
// Using the authenticated client
const client = await createAuthenticatedClient();
const profile = await client.get('/profile');
```

Or manually:
```javascript
const response = await fetch('/api/v1/profile', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'X-Platform-ID': 'vk',
    'X-Child-Profile-ID': sessionStorage.getItem('child_profile_id')
  }
});
```

## Testing

### Manual Testing Steps
1. Open VK Mini App
2. Check browser console for:
   - `[Auth] Using cached profile ID: <uuid>` or
   - `[Auth] Profile found: <uuid>`
3. Verify profile ID in sessionStorage
4. Make API requests and check Network tab for:
   - `X-Platform-ID` header present
   - `X-Child-Profile-ID` header present
   - 200 OK responses

### Automated Testing
Run the backend tests:
```bash
cd api
go test ./... -v
```

Test CORS configuration:
```bash
curl -H "Origin: http://77.222.60.149" \
     -H "Access-Control-Request-Method: GET" \
     -X OPTIONS http://localhost:8080/api/v1/health
```

Expected: `Access-Control-Allow-Origin: http://77.222.60.149`

## Security Considerations

### Session Storage vs Local Storage
- Profile ID stored in `sessionStorage` (cleared on tab close)
- More secure than `localStorage` for sensitive data
- Automatically cleared when browser/tab closes

### CORS Configuration
- Production: Strict origin checking
- Development/Staging: Permissive for developer convenience
- Credentials included only for allowed origins

### VK Signature Validation
- Backend validates `sign` parameter using VK App secret
- Prevents request forgery
- Required for all VK-originated requests in production

## Troubleshooting

### Child Profile ID Still Null
1. Check browser console for errors
2. Verify `sessionStorage` contains `child_profile_id`
3. Ensure VK Bridge initialized successfully
4. Check network tab for `/profiles/by-platform` response

### CORS Errors
1. Verify `ALLOWED_ORIGINS` in `.env.production`
2. Check request includes correct `Origin` header
3. Ensure no protocol mismatch (HTTP vs HTTPS)
4. Check nginx/AWS load balancer CORS settings

### API Returns 401 Unauthorized
1. Verify `X-Platform-ID` header present
2. Verify `X-Child-Profile-ID` header present and valid
3. Check auth middleware logs
4. Ensure profile exists in database

## Migration Notes

### For Existing Users
- Existing profiles remain valid
- New auth flow backwards compatible
- Cached profile IDs reused automatically

### For New Users
- First visit triggers profile creation
- Profile ID cached immediately
- Subsequent visits use cached ID

## Performance Impact

- **Positive**: Reduced API calls (cached profile ID)
- **Positive**: Faster authorization (no database lookup for cached IDs)
- **Neutral**: Minimal overhead for header injection
- **Positive**: Better error handling and logging

## Future Improvements

1. Implement JWT tokens for stateless authentication
2. Add refresh token mechanism
3. Implement rate limiting per profile
4. Add audit logging for auth events
5. Implement multi-device sync for profile data
EOF