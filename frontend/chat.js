class ChatApp {
    constructor() {
        this.messages = [];
        this.isLoading = false;
        
        this.messagesContainer = document.getElementById('messages');
        this.messageInput = document.getElementById('messageInput');
        this.chatButton = document.getElementById('chatButton');
        this.chatWithToolButton = document.getElementById('chatWithToolButton');
        this.chatWithAgentButton = document.getElementById('chatWithAgentButton');
        this.chatWithDocButton = document.getElementById('chatWithDocButton');
        
        this.init();
    }
    
    init() {
        this.chatButton.addEventListener('click', () => this.sendMessage('chat'));
        this.chatWithToolButton.addEventListener('click', () => this.sendMessage('tool'));
        this.chatWithAgentButton.addEventListener('click', () => this.sendMessage('agent'));
        this.chatWithDocButton.addEventListener('click', () => this.sendMessage('doc'));
        
        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage('chat'); // Default to chat on Enter
            }
        });
    }
    
    async sendMessage(mode = 'chat') {
        const content = this.messageInput.value.trim();
        if (!content || this.isLoading) return;
        
        // Add user message
        this.addMessage('user', content);
        this.messageInput.value = '';
        
        // Show typing indicator
        this.showTypingIndicator(mode);
        
        try {
            let response;
            switch (mode) {
                case 'chat':
                    response = await this.callChatAPI(content);
                    break;
                case 'tool':
                    response = await this.callChatWithToolAPI(content);
                    break;
                case 'agent':
                    response = await this.callChatWithAgentAPI(content);
                    break;
                case 'doc':
                    response = await this.callChatWithDocAPI(content);
                    break;
                default:
                    response = await this.callChatAPI(content);
            }
            
            // Hide typing indicator
            this.hideTypingIndicator();
            
            // Add assistant response
            this.addMessage('assistant', response.content);
            
        } catch (error) {
            this.hideTypingIndicator();
            this.showError(`Failed to get ${mode} response: ` + error.message);
        }
    }
    
    async callChatAPI(userMessage) {
        // Prepare messages for the API
        const messages = [
            ...this.messages,
            { role: 'ROLE_USER', content: userMessage }
        ];
        
        const requestData = {
            messages: messages,
            temperature: 0.7,
            max_tokens: 500
        };
        
        // Since this is a demo and we can't make direct gRPC calls from browser,
        // we'll simulate the API call. In a real implementation, you would have
        // a REST API gateway that translates HTTP to gRPC.
        return this.simulateAPICall(requestData, 'chat');
    }
    
    async callChatWithToolAPI(userMessage) {
        const messages = [
            ...this.messages,
            { role: 'ROLE_USER', content: userMessage }
        ];
        
        const requestData = {
            messages: messages,
            temperature: 0.7,
            max_tokens: 500
        };
        
        return this.simulateAPICall(requestData, 'tool');
    }
    
    async callChatWithAgentAPI(userMessage) {
        const messages = [
            ...this.messages,
            { role: 'ROLE_USER', content: userMessage }
        ];
        
        const requestData = {
            messages: messages,
            temperature: 0.7,
            max_tokens: 500
        };
        
        return this.simulateAPICall(requestData, 'agent');
    }
    
    async callChatWithDocAPI(userMessage) {
        const messages = [
            ...this.messages,
            { role: 'ROLE_USER', content: userMessage }
        ];
        
        const requestData = {
            messages: messages,
            temperature: 0.7,
            max_tokens: 500
        };
        
        return this.simulateAPICall(requestData, 'doc');
    }
    
    async simulateAPICall(requestData, mode = 'chat') {
        // Simulate network delay
        await new Promise(resolve => setTimeout(resolve, 1000 + Math.random() * 2000));
        
        // Simulate different types of responses based on input and mode
        const userMessage = requestData.messages[requestData.messages.length - 1].content.toLowerCase();
        
        let responseContent;
        const modePrefix = mode === 'chat' ? '' : `[${mode.charAt(0).toUpperCase() + mode.slice(1)} Mode] `;
        
        if (userMessage.includes('hello') || userMessage.includes('hi')) {
            responseContent = modePrefix + "Hello! I'm your AI assistant. How can I help you today?";
        } else if (userMessage.includes('weather')) {
            if (mode === 'tool') {
                responseContent = modePrefix + "🔍 Using weather tools... The current weather is partly cloudy with a temperature of 22°C. Tool-enhanced response with real-time data access!";
            } else if (mode === 'agent') {
                responseContent = modePrefix + "🤖 As an intelligent agent, I can coordinate multiple services to get weather data, set reminders, and suggest activities based on conditions.";
            } else if (mode === 'doc') {
                responseContent = modePrefix + "📄 Based on weather documentation and patterns, here's comprehensive weather analysis with historical comparisons and forecasts.";
            } else {
                responseContent = "I don't have access to real-time weather data, but I'd be happy to help you with other questions!";
            }
        } else if (userMessage.includes('time') || userMessage.includes('date')) {
            responseContent = modePrefix + `The current time is ${new Date().toLocaleString()}. Is there anything else I can help you with?`;
        } else if (userMessage.includes('joke')) {
            const jokes = [
                "Why don't scientists trust atoms? Because they make up everything!",
                "Why did the scarecrow win an award? He was outstanding in his field!",
                "Why don't eggs tell jokes? They'd crack each other up!"
            ];
            responseContent = modePrefix + jokes[Math.floor(Math.random() * jokes.length)];
        } else if (userMessage.includes('help')) {
            let helpText = "I'm here to help! You can ask me questions about various topics, request jokes, ask about the time, or just have a conversation.";
            if (mode === 'tool') {
                helpText += " In Tool mode, I can search the web, use calculators, and access external APIs.";
            } else if (mode === 'agent') {
                helpText += " In Agent mode, I can coordinate multiple tasks, plan workflows, and act autonomously.";
            } else if (mode === 'doc') {
                helpText += " In Doc mode, I can analyze documents, extract information, and provide detailed research.";
            }
            responseContent = modePrefix + helpText;
        } else {
            let baseResponse = `You said: "${requestData.messages[requestData.messages.length - 1].content}". That's interesting!`;
            if (mode === 'tool') {
                baseResponse += " I can enhance this with web search, calculations, and external data.";
            } else if (mode === 'agent') {
                baseResponse += " As an agent, I can break this down into actionable steps and coordinate resources.";
            } else if (mode === 'doc') {
                baseResponse += " I can provide detailed analysis and cross-reference with relevant documentation.";
            } else {
                baseResponse += " I'm a demo AI assistant powered by a Go gRPC backend. Feel free to ask me anything!";
            }
            responseContent = modePrefix + baseResponse;
        }
        
        return {
            content: responseContent,
            token_usage: {
                input_token_num: Math.floor(Math.random() * 50) + 10,
                output_token_num: Math.floor(Math.random() * 30) + 10,
                total_token_num: Math.floor(Math.random() * 80) + 20
            }
        };
    }
    
    
    addMessage(role, content) {
        // Add to messages array
        const roleMap = {
            'user': 'ROLE_USER',
            'assistant': 'ROLE_ASSISTANT',
            'system': 'ROLE_SYSTEM'
        };
        
        this.messages.push({
            role: roleMap[role] || 'ROLE_USER',
            content: content
        });
        
        // Add to UI
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${role}`;
        messageDiv.textContent = content;
        
        this.messagesContainer.appendChild(messageDiv);
        this.scrollToBottom();
    }
    
    showTypingIndicator(mode = 'chat') {
        this.isLoading = true;
        
        // Disable all buttons
        this.chatButton.disabled = true;
        this.chatWithToolButton.disabled = true;
        this.chatWithAgentButton.disabled = true;
        this.chatWithDocButton.disabled = true;
        
        // Update button text to show processing state
        this.chatButton.textContent = 'Processing...';
        this.chatWithToolButton.textContent = 'Processing...';
        this.chatWithAgentButton.textContent = 'Processing...';
        this.chatWithDocButton.textContent = 'Processing...';
        
        const typingDiv = document.createElement('div');
        typingDiv.className = 'typing-indicator';
        typingDiv.id = 'typing-indicator';
        
        let modeText = mode.charAt(0).toUpperCase() + mode.slice(1);
        typingDiv.innerHTML = `${modeText} AI is typing<span class="typing-dots"></span>`;
        
        this.messagesContainer.appendChild(typingDiv);
        this.scrollToBottom();
    }
    
    hideTypingIndicator() {
        this.isLoading = false;
        
        // Enable all buttons
        this.chatButton.disabled = false;
        this.chatWithToolButton.disabled = false;
        this.chatWithAgentButton.disabled = false;
        this.chatWithDocButton.disabled = false;
        
        // Restore button text
        this.chatButton.textContent = 'Chat';
        this.chatWithToolButton.textContent = 'Tool';
        this.chatWithAgentButton.textContent = 'Agent';
        this.chatWithDocButton.textContent = 'Doc';
        
        const typingIndicator = document.getElementById('typing-indicator');
        if (typingIndicator) {
            typingIndicator.remove();
        }
    }
    
    showError(message) {
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error';
        errorDiv.textContent = message;
        
        this.messagesContainer.appendChild(errorDiv);
        this.scrollToBottom();
        
        // Remove error after 5 seconds
        setTimeout(() => {
            if (errorDiv.parentNode) {
                errorDiv.remove();
            }
        }, 5000);
    }
    
    scrollToBottom() {
        this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
    }
}

// Initialize the chat app when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new ChatApp();
});