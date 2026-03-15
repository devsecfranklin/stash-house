// This is highly simplified and assumes you have a way to receive chat messages
// For a real application, you'd use a WebSocket connection to a server that
// fetches messages from Twitch's API (EventSub or IRC).

const chatContainer = document.getElementById('chat-container');
const MAX_MESSAGES = 20; // Maximum number of messages to display
const MESSAGE_LIFETIME = 15000; // Messages disappear after 15 seconds (in ms)

function addChatMessage(username, message) {
    const messageElement = document.createElement('div');
    messageElement.classList.add('chat-message');

    const usernameSpan = document.createElement('span');
    usernameSpan.classList.add('username');
    usernameSpan.textContent = username + ':';

    const messageTextSpan = document.createElement('span');
    messageTextSpan.classList.add('message-text');
    messageTextSpan.textContent = message;

    messageElement.appendChild(usernameSpan);
    messageElement.appendChild(messageTextSpan);

    chatContainer.prepend(messageElement); // Add to the top for flex-direction: column-reverse

    // Remove old messages if too many
    while (chatContainer.children.length > MAX_MESSAGES) {
        chatContainer.lastChild.remove();
    }

    // Schedule message to fade out and be removed
    setTimeout(() => {
        messageElement.classList.add('fade-out');
        messageElement.addEventListener('animationend', () => {
            messageElement.remove();
        }, { once: true });
    }, MESSAGE_LIFETIME);
}

// --- DEMO FUNCTIONALITY (REMOVE FOR PRODUCTION) ---
// Simulate receiving messages
let messageCount = 0;
setInterval(() => {
    messageCount++;
    const usernames = ["StreamerBot", "Viewer123", "EpicGamer", "ChattyPatty", "Moderator"];
    const messages = [
        "Hello everyone, how's it going?",
        "This game is so much fun!",
        " poggers",
        "Can you play this song next?",
        "Thanks for the raid!",
        "What's your favorite part of streaming?",
        "Wow, nice play!",
        "LUL",
        "Kappa",
        "This is a longer message to test wrapping capabilities and see how it handles multiple lines of text.",
        "Another test message from another viewer."
    ];
    const randomUsername = usernames[Math.floor(Math.random() * usernames.length)];
    const randomMessage = messages[Math.floor(Math.random() * messages.length)];
    addChatMessage(randomUsername, randomMessage);
}, 2000); // Add a new message every 2 seconds
// --- END DEMO FUNCTIONALITY ---

// In a real scenario, you'd have something like:
/*
const ws = new WebSocket('ws://localhost:8080'); // Your WebSocket server

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.type === 'chatMessage') {
        addChatMessage(data.username, data.message);
    }
};

ws.onopen = () => {
    console.log('Connected to WebSocket server');
};

ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};

ws.onclose = () => {
    console.log('WebSocket connection closed');
};
*/
