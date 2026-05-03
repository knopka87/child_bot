// Тестовая страница для проверки параметров URL
export function TestParams() {
  const urlParams = new URLSearchParams(window.location.search);
  const hashParams = new URLSearchParams(window.location.hash.substring(1));

  const allParams: Record<string, string> = {};

  // Собираем все параметры из query
  urlParams.forEach((value, key) => {
    allParams[`query:${key}`] = value;
  });

  // Собираем все параметры из hash
  hashParams.forEach((value, key) => {
    allParams[`hash:${key}`] = value;
  });

  return (
    <div style={{ padding: '20px', fontFamily: 'monospace', fontSize: '12px' }}>
      <h1>URL Parameters Debug</h1>

      <h2>Full URL:</h2>
      <pre style={{ background: '#f0f0f0', padding: '10px', overflow: 'auto' }}>
        {window.location.href}
      </pre>

      <h2>window.location.search:</h2>
      <pre style={{ background: '#f0f0f0', padding: '10px' }}>
        {window.location.search || '(empty)'}
      </pre>

      <h2>window.location.hash:</h2>
      <pre style={{ background: '#f0f0f0', padding: '10px' }}>
        {window.location.hash || '(empty)'}
      </pre>

      <h2>All Parsed Parameters:</h2>
      <pre style={{ background: '#f0f0f0', padding: '10px' }}>
        {Object.keys(allParams).length > 0
          ? JSON.stringify(allParams, null, 2)
          : '(no parameters found)'}
      </pre>

      <h2>Looking for vk_ref:</h2>
      <pre style={{ background: '#f0f0f0', padding: '10px' }}>
        {`query vk_ref: ${urlParams.get('vk_ref') || 'NOT FOUND'}\n`}
        {`hash vk_ref: ${hashParams.get('vk_ref') || 'NOT FOUND'}`}
      </pre>
    </div>
  );
}
