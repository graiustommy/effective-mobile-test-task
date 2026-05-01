CREATE TABLE IF NOT EXISTS subscriptions(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid()
	service_name varchar(255) UNIQUE,
	price int,
	user_id UUID NOT NULL UNIQUE, 
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
);

CREATE INDEX user_id_index ON subscriptions(user_id);
CREATE INDEX service_name_index ON subscriptions(service_name);
