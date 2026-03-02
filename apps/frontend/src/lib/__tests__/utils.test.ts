import { describe, it, expect } from 'vitest';

describe('Utility Functions', () => {
  describe('String Utils', () => {
    it('should truncate string correctly', () => {
      const truncate = (str: string, length: number): string => {
        if (str.length <= length) return str;
        return str.slice(0, length) + '...';
      };

      expect(truncate('Hello World', 5)).toBe('Hello...');
      expect(truncate('Hi', 5)).toBe('Hi');
    });
  });

  describe('Number Utils', () => {
    it('should format number with commas', () => {
      const formatNumber = (num: number): string => {
        return num.toLocaleString('zh-CN');
      };

      expect(formatNumber(1000)).toBe('1,000');
      expect(formatNumber(1234567)).toBe('1,234,567');
    });
  });

  describe('Date Utils', () => {
    it('should format date', () => {
      const formatDate = (date: Date): string => {
        return date.toISOString().split('T')[0];
      };

      const date = new Date('2024-02-28T10:30:00Z');
      expect(formatDate(date)).toBe('2024-02-28');
    });
  });

  describe('Validation', () => {
    it('should validate email format', () => {
      const isValidEmail = (email: string): boolean => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
      };

      expect(isValidEmail('test@example.com')).toBe(true);
      expect(isValidEmail('invalid-email')).toBe(false);
    });
  });
});
