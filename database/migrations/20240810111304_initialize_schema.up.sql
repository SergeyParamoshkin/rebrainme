-- Создание таблицы alerts
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы tickets
CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    queue VARCHAR(255),
    key VARCHAR(255) UNIQUE NOT NULL,
    summary VARCHAR(255) NOT NULL,
    type VARCHAR(255)
);

-- Создание таблицы связей alert_ticket
CREATE TABLE IF NOT EXISTS alert_ticket (
    alert_id INT,
    ticket_id UUID,
    PRIMARY KEY (alert_id, ticket_id),
    CONSTRAINT fk_alert FOREIGN KEY (alert_id) REFERENCES alerts (id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket FOREIGN KEY (ticket_id) REFERENCES tickets (id) ON DELETE CASCADE
);

-- Индексы для таблицы alert_ticket
CREATE INDEX idx_alert_ticket_alert_id ON alert_ticket (alert_id);

CREATE INDEX idx_alert_ticket_ticket_id ON alert_ticket (ticket_id);