CREATE TABLE subscriptions (
                               id SERIAL PRIMARY KEY,
                               user_id VARCHAR(36) NOT NULL,
                               service_name VARCHAR(255) NOT NULL,
                               price INTEGER NOT NULL,
                               start_date DATE NOT NULL,
                               end_date DATE
);

CREATE INDEX idx_subscriptions_user_service_dates
    ON subscriptions (user_id, service_name, start_date, end_date);