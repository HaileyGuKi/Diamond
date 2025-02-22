const puppeteer = require('puppeteer');
const fs = require('fs');
const FormData = require('form-data');
const axios = require('axios');
require('dotenv').config();

(async () => {
    try {
        const logFilePath = '../liga-table.log'; // Update the path if necessary

        // Step 1: Check the log file for LEAGUE
        const logData = fs.readFileSync(logFilePath, 'utf-8');
        const hasLeague = logData.split('\n').some(line => line.includes('LEAGUE'));
        if (!hasLeague) {
            console.log('No LEAGUE entry found in liga-table.log. Exiting script.');
            return;
        }

        // Step 2: Launch Puppeteer to scrape the table
        const browser = await puppeteer.launch({
            args: ['--no-sandbox', '--disable-setuid-sandbox'], // Add these flags
        });
        const page = await browser.newPage();

        // Set a high-resolution viewport
        await page.setViewport({
            width: 1920,
            height: 1080,
            deviceScaleFactor: 2,
        });

        // Navigate to the La Liga table page
        const url = 'https://www.bbc.com/sport/football/spanish-la-liga/table';
        await page.goto(url, { waitUntil: 'networkidle2' });

        // Dismiss the cookie banner if it exists
        const cookieAcceptButtonSelector = 'button[data-testid="accept-button"]';
        try {
            await page.waitForSelector(cookieAcceptButtonSelector, { timeout: 5000 });
            await page.click(cookieAcceptButtonSelector);
            console.log('Cookie banner dismissed by accepting cookies.');
        } catch (e) {
            console.log('No cookie banner found or unable to dismiss.');
        }

        // Wait for the table element
        const tableSelector = '.ssrcss-1dbj4ao-TableWrapper';
        await page.waitForSelector(tableSelector);

        // Take a screenshot of the table element
        const imagePath = 'liga-high-quality.png';
        const tableElement = await page.$(tableSelector);
        if (tableElement) {
            await tableElement.screenshot({ path: imagePath });
            console.log(`High-quality screenshot saved as ${imagePath}`);
        } else {
            console.log('Table element not found!');
        }

        await browser.close();

        // Step 3: Send the screenshot to Telegram
        if (fs.existsSync(imagePath)) {
            await sendImageToTelegram(imagePath);
            console.log('Image sent to Telegram successfully.');

            // Remove the image file
            fs.unlinkSync(imagePath);
            console.log(`${imagePath} removed from the project folder.`);
        } else {
            console.error('Image file not found for sending to Telegram.');
        }

        // Step 4: Clear the log file
        fs.writeFileSync(logFilePath, '');
        console.log(`${logFilePath} has been cleared.`);
    } catch (error) {
        console.error('Error:', error);
    }
})();

async function sendImageToTelegram(imagePath) {
    const TELEGRAM_BOT_TOKEN = process.env.TELEGRAM_BOT_TOKEN;
    const TELEGRAM_CHANNEL_ID = process.env.TELEGRAM_CHANNEL_ID;

    if (!TELEGRAM_BOT_TOKEN || !TELEGRAM_CHANNEL_ID) {
        throw new Error('Missing Telegram bot token or channel ID in .env file.');
    }

    // Create a form and attach the file
    const formData = new FormData();
    formData.append('chat_id', TELEGRAM_CHANNEL_ID);
    formData.append('photo', fs.createReadStream(imagePath));

    try {
        // Send the form data to Telegram
        const response = await axios.post(
            `https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendPhoto`,
            formData,
            { headers: formData.getHeaders() }
        );

        if (response.data && response.data.ok) {
            console.log('Image successfully sent to Telegram!');
        } else {
            console.error('Failed to send image to Telegram:', response.data);
        }
    } catch (error) {
        console.error('Error sending image to Telegram:', error.response?.data || error.message);
    }
}
