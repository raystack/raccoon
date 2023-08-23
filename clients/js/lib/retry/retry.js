async function retry(callback, maxAttempts, waitTime, logger) {
    for (let attempt = 0; attempt < maxAttempts; attempt += 1) {
        try {
            // eslint-disable-next-line no-await-in-loop
            return await callback();
        } catch (error) {
            if (attempt === maxAttempts - 1) {
                throw error;
            }
            // eslint-disable-next-line
            await new Promise((resolve) => setTimeout(resolve, waitTime));
            logger.info(
                `[Retry ${attempt}]: Retrying after ${waitTime} due to the Error: ${error}`
            );
        }
    }
    throw new Error('Retry limit exceeded');
}

export default retry;
