// Пример клиентского кода для Telegram Mini App
class TelegramMiniAppClient {
    constructor(baseUrl) {
        this.baseUrl = baseUrl;
        this.token = null;
        this.user = null;
    }

    // Инициализация с Telegram init data
    async init(initData) {
        try {
            const response = await fetch(`${this.baseUrl}/api/v1/auth/telegram`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ init_data: initData })
            });

            if (!response.ok) {
                throw new Error('Authentication failed');
            }

            const data = await response.json();
            this.token = data.token;
            this.user = data.user;

            // Сохраняем токен в localStorage
            localStorage.setItem('tma_token', this.token);
            localStorage.setItem('tma_user', JSON.stringify(this.user));

            return data;
        } catch (error) {
            console.error('Authentication error:', error);
            throw error;
        }
    }

    // Восстановление сессии из localStorage
    restoreSession() {
        const token = localStorage.getItem('tma_token');
        const user = localStorage.getItem('tma_user');

        if (token && user) {
            this.token = token;
            this.user = JSON.parse(user);
            return true;
        }
        return false;
    }

    // Выход из системы
    logout() {
        this.token = null;
        this.user = null;
        localStorage.removeItem('tma_token');
        localStorage.removeItem('tma_user');
    }

    // Получение заголовков для авторизованных запросов
    getAuthHeaders() {
        if (!this.token) {
            throw new Error('Not authenticated');
        }
        return {
            'Authorization': `Bearer ${this.token}`,
            'Content-Type': 'application/json',
        };
    }

    // Получение профиля пользователя
    async getProfile() {
        const response = await fetch(`${this.baseUrl}/api/v1/user/profile`, {
            headers: this.getAuthHeaders()
        });

        if (!response.ok) {
            throw new Error('Failed to get profile');
        }

        return await response.json();
    }

    // Получение всех элементов
    async getItems() {
        const response = await fetch(`${this.baseUrl}/api/v1/items`, {
            headers: this.getAuthHeaders()
        });

        if (!response.ok) {
            throw new Error('Failed to get items');
        }

        return await response.json();
    }

    // Получение элемента по ID
    async getItem(id) {
        const response = await fetch(`${this.baseUrl}/api/v1/items/${id}`, {
            headers: this.getAuthHeaders()
        });

        if (!response.ok) {
            throw new Error('Failed to get item');
        }

        return await response.json();
    }

    // Создание нового элемента
    async createItem(title, description = '') {
        const response = await fetch(`${this.baseUrl}/api/v1/items`, {
            method: 'POST',
            headers: this.getAuthHeaders(),
            body: JSON.stringify({ title, description })
        });

        if (!response.ok) {
            throw new Error('Failed to create item');
        }

        return await response.json();
    }

    // Обновление элемента
    async updateItem(id, updates) {
        const response = await fetch(`${this.baseUrl}/api/v1/items/${id}`, {
            method: 'PUT',
            headers: this.getAuthHeaders(),
            body: JSON.stringify(updates)
        });

        if (!response.ok) {
            throw new Error('Failed to update item');
        }

        return await response.json();
    }

    // Удаление элемента
    async deleteItem(id) {
        const response = await fetch(`${this.baseUrl}/api/v1/items/${id}`, {
            method: 'DELETE',
            headers: this.getAuthHeaders()
        });

        if (!response.ok) {
            throw new Error('Failed to delete item');
        }

        return await response.json();
    }
}

// Пример использования
async function example() {
    // Инициализация клиента
    const client = new TelegramMiniAppClient('https://your-app.railway.app');

    // Попытка восстановить сессию
    if (!client.restoreSession()) {
        // Если нет сохраненной сессии, инициализируем новую
        const initData = window.Telegram.WebApp.initData;
        await client.init(initData);
    }

    try {
        // Получаем профиль пользователя
        const profile = await client.getProfile();
        console.log('User profile:', profile);

        // Получаем все элементы
        const items = await client.getItems();
        console.log('Items:', items);

        // Создаем новый элемент
        const newItem = await client.createItem('Новая задача', 'Описание задачи');
        console.log('Created item:', newItem);

        // Обновляем элемент
        const updatedItem = await client.updateItem(newItem.id, {
            title: 'Обновленная задача',
            status: 'completed'
        });
        console.log('Updated item:', updatedItem);

    } catch (error) {
        console.error('API error:', error);
    }
}

// Экспорт для использования в других модулях
if (typeof module !== 'undefined' && module.exports) {
    module.exports = TelegramMiniAppClient;
}
