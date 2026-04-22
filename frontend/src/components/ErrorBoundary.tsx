// src/components/ErrorBoundary.tsx
import { Component, ErrorInfo, ReactNode } from 'react';
import { Div, Title, Text, Button } from '@vkontakte/vkui';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
}

/**
 * Error Boundary для отлова ошибок React
 */
export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
    };
  }

  static getDerivedStateFromError(error: Error): Partial<State> {
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('[ErrorBoundary] Caught error:', error, errorInfo);

    this.setState({
      error,
      errorInfo,
    });

    // Здесь можно отправить ошибку в систему мониторинга
    // например, Sentry
  }

  handleReload = () => {
    window.location.reload();
  };

  render() {
    if (this.state.hasError) {
      return (
        <Div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: '100vh',
            padding: '20px',
            textAlign: 'center',
          }}
        >
          <Title level="1" weight="1" style={{ marginBottom: '16px' }}>
            Что-то пошло не так 😔
          </Title>
          <Text style={{ marginBottom: '24px', color: 'var(--vkui--color_text_secondary)' }}>
            Произошла ошибка в приложении. Попробуйте перезагрузить страницу.
          </Text>

          {import.meta.env.DEV && this.state.error && (
            <div
              style={{
                marginBottom: '24px',
                padding: '16px',
                background: 'var(--vkui--color_background_secondary)',
                borderRadius: '8px',
                textAlign: 'left',
                maxWidth: '100%',
                overflow: 'auto',
              }}
            >
              <Text weight="2" style={{ marginBottom: '8px' }}>
                Error:
              </Text>
              <Text style={{ fontFamily: 'monospace', fontSize: '12px' }}>
                {this.state.error.toString()}
              </Text>
              {this.state.errorInfo && (
                <>
                  <Text weight="2" style={{ marginTop: '12px', marginBottom: '8px' }}>
                    Stack trace:
                  </Text>
                  <Text
                    style={{
                      fontFamily: 'monospace',
                      fontSize: '10px',
                      whiteSpace: 'pre-wrap',
                    }}
                  >
                    {this.state.errorInfo.componentStack}
                  </Text>
                </>
              )}
            </div>
          )}

          <Button size="l" mode="primary" onClick={this.handleReload} stretched>
            Перезагрузить
          </Button>
        </Div>
      );
    }

    return this.props.children;
  }
}
