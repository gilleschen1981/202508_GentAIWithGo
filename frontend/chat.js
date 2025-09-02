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
        
        return this.callHTTPAPI('/api/chat', requestData);
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
        
        return this.callHTTPAPI('/api/chat-with-tool', requestData);
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
        
        return this.callHTTPAPI('/api/chat-with-agent', requestData);
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
        
        return this.callHTTPAPI('/api/chat-with-doc', requestData);
    }
    
    async callHTTPAPI(endpoint, requestData) {
        const baseURL = 'http://localhost:8080';
        const url = baseURL + endpoint;
        
        console.log(`🌐 Calling API: ${url}`);
        console.log('📤 Request data:', requestData);
        
        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });
            
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP ${response.status}: ${errorText}`);
            }
            
            const responseData = await response.json();
            console.log('📥 Response data:', responseData);
            
            if (responseData.error) {
                throw new Error(responseData.error);
            }
            
            return responseData;
            
        } catch (error) {
            console.error('❌ API call failed:', error);
            throw new Error(`Failed to connect to server: ${error.message}`);
        }
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