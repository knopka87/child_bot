// src/lib/logger.ts
/**
 * Структурированный logger с уровнями логирования
 * Заменяет прямое использование console.log/error/warn
 */

import config from '@/config';

export enum LogLevel {
  DEBUG = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3,
  NONE = 4,
}

export interface LogContext {
  module?: string;
  userId?: string;
  sessionId?: string;
  childProfileId?: string;
  [key: string]: any;
}

interface LogEntry {
  timestamp: string;
  level: string;
  message: string;
  context: LogContext;
  data?: any;
}

export class Logger {
  private level: LogLevel;
  private context: LogContext = {};

  constructor(level?: LogLevel) {
    // В production только WARN и ERROR
    // В development все логи
    this.level =
      level !== undefined
        ? level
        : config.isProduction
        ? LogLevel.WARN
        : LogLevel.DEBUG;
  }

  /**
   * Устанавливает контекст для всех последующих логов
   */
  setContext(context: LogContext): void {
    this.context = { ...this.context, ...context };
  }

  /**
   * Очищает контекст
   */
  clearContext(): void {
    this.context = {};
  }

  /**
   * Обновляет контекст (merge)
   */
  updateContext(context: Partial<LogContext>): void {
    this.context = { ...this.context, ...context };
  }

  private log(level: LogLevel, message: string, data?: any): void {
    if (level < this.level) return;

    const timestamp = new Date().toISOString();
    const logEntry: LogEntry = {
      timestamp,
      level: LogLevel[level],
      message,
      context: this.context,
      ...(data !== undefined && { data }),
    };

    // В production отправляем ERROR в error tracking service
    if (config.isProduction && level >= LogLevel.ERROR) {
      this.sendToErrorTracking(logEntry);
    }

    // Console output с форматированием
    this.outputToConsole(level, message, data);
  }

  private outputToConsole(level: LogLevel, message: string, data?: any): void {
    const prefix = this.context.module
      ? `[${this.context.module}]`
      : '[App]';

    const formattedMessage = `${prefix} ${message}`;

    switch (level) {
      case LogLevel.DEBUG:
        // eslint-disable-next-line no-console
        console.debug(formattedMessage, data !== undefined ? data : '');
        break;
      case LogLevel.INFO:
        // eslint-disable-next-line no-console
        console.info(formattedMessage, data !== undefined ? data : '');
        break;
      case LogLevel.WARN:
        // eslint-disable-next-line no-console
        console.warn(formattedMessage, data !== undefined ? data : '');
        break;
      case LogLevel.ERROR:
        // eslint-disable-next-line no-console
        console.error(formattedMessage, data !== undefined ? data : '');
        break;
    }
  }

  private sendToErrorTracking(logEntry: LogEntry): void {
    // TODO: Интеграция с Sentry или другим error tracking service
    // Пример:
    // if (window.Sentry) {
    //   window.Sentry.captureException(new Error(logEntry.message), {
    //     level: logEntry.level.toLowerCase() as any,
    //     extra: {
    //       ...logEntry.context,
    //       ...logEntry.data,
    //     },
    //   });
    // }

    // Пока просто сохраняем в localStorage для отладки
    try {
      const errors = JSON.parse(localStorage.getItem('app_errors') || '[]');
      errors.push({
        ...logEntry,
        userAgent: navigator.userAgent,
        url: window.location.href,
      });

      // Храним максимум 50 последних ошибок
      if (errors.length > 50) {
        errors.shift();
      }

      localStorage.setItem('app_errors', JSON.stringify(errors));
    } catch (e) {
      // Игнорируем ошибки при сохранении в localStorage
    }
  }

  /**
   * Debug level - подробная информация для разработки
   */
  debug(message: string, data?: any): void {
    this.log(LogLevel.DEBUG, message, data);
  }

  /**
   * Info level - общая информация о работе приложения
   */
  info(message: string, data?: any): void {
    this.log(LogLevel.INFO, message, data);
  }

  /**
   * Warn level - предупреждения о потенциальных проблемах
   */
  warn(message: string, data?: any): void {
    this.log(LogLevel.WARN, message, data);
  }

  /**
   * Error level - критические ошибки
   */
  error(message: string, data?: any): void {
    this.log(LogLevel.ERROR, message, data);
  }

  /**
   * Создаёт группу логов (для визуальной группировки в консоли)
   */
  group(label: string): void {
    if (this.level > LogLevel.DEBUG) return;
    // eslint-disable-next-line no-console
    console.group(`[${this.context.module || 'App'}] ${label}`);
  }

  /**
   * Закрывает группу логов
   */
  groupEnd(): void {
    if (this.level > LogLevel.DEBUG) return;
    // eslint-disable-next-line no-console
    console.groupEnd();
  }

  /**
   * Логирует таблицу (для массивов объектов)
   */
  table(data: any[]): void {
    if (this.level > LogLevel.DEBUG) return;
    // eslint-disable-next-line no-console
    console.table(data);
  }
}

// Глобальный инстанс логгера
export const logger = new Logger();

/**
 * Создаёт logger с контекстом модуля
 */
export function createLogger(module: string, context?: Omit<LogContext, 'module'>): Logger {
  const moduleLogger = new Logger();
  moduleLogger.setContext({ module, ...context });
  return moduleLogger;
}

/**
 * Helper для логирования performance metrics
 */
export function logPerformance(name: string, duration: number, context?: Record<string, any>): void {
  if (config.isDevelopment) {
    logger.debug(`Performance: ${name}`, {
      duration_ms: Math.round(duration),
      ...context,
    });
  }
}
