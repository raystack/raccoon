// eslint-disable-next-line
import { jest } from '@jest/globals';
import retry from '../../lib/retry/retry.js';

describe('retry', () => {
    test('should return callback result on first success attempt', async () => {
        const mockCallback = jest.fn(() => 'success');
        const result = await retry(mockCallback, 3, 100, console);

        expect(result).toEqual('success');
        expect(mockCallback).toHaveBeenCalledTimes(1);
    });

    test('should return callback result after multiple attempts', async () => {
        const mockCallback = jest
            .fn()
            .mockRejectedValueOnce(new Error('Attempt 1 failed'))
            .mockRejectedValueOnce(new Error('Attempt 2 failed'))
            .mockResolvedValueOnce('success');

        const result = await retry(mockCallback, 3, 100, console);

        expect(result).toEqual('success');
        expect(mockCallback).toHaveBeenCalledTimes(3);
    });

    test('should throw error if callback always fails', async () => {
        const mockCallback = jest.fn().mockRejectedValue(new Error('All attempts failed'));

        await expect(retry(mockCallback, 3, 100, console)).rejects.toThrow('All attempts failed');
        expect(mockCallback).toHaveBeenCalledTimes(3);
    });

    test('should throw error if retry limit exceeded', async () => {
        const mockCallback = jest
            .fn()
            .mockRejectedValueOnce(new Error('Attempt 1 failed'))
            .mockRejectedValueOnce(new Error('Attempt 2 failed'))
            .mockRejectedValueOnce(new Error('Attempt 3 failed'))
            .mockResolvedValueOnce('success');

        await expect(retry(mockCallback, 3, 100, console)).rejects.toThrow('Attempt 3 failed');
        expect(mockCallback).toHaveBeenCalledTimes(3);
    });

    test('should apply waitTime', async () => {
        const startTime = Date.now();
        const mockCallback = jest.fn().mockRejectedValue(new Error('Attempt failed'));

        try {
            await retry(mockCallback, 3, 500, console);
        } catch (error) {
            const endTime = Date.now();
            const elapsedTime = endTime - startTime;
            expect(elapsedTime).toBeGreaterThanOrEqual(1000);
            expect(elapsedTime).toBeLessThanOrEqual(1500);
        }

        expect(mockCallback).toHaveBeenCalledTimes(3);
    });
});
