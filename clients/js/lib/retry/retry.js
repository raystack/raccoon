export async function retry(callback, maxAttempts, waitTime) {
    for (let attempt = 0; attempt < maxAttempts; attempt++) {
        try {
            return await callback();
        } catch (error) {
            if (attempt === maxAttempts - 1) {
                throw error;
            }
            await new Promise(resolve => setTimeout(resolve, waitTime));
            console.info(`[Retry ${attempt}]: Retrying after ${waitTime} due to the Error: ${error}`);
        }
    }
    throw new Error('Retry limit exceeded');
}
