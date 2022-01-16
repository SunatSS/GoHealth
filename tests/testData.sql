INSERT INTO customers (name, phone, password, address) VALUES ('Fedya', '1', '12345678', '123 Main St') ON CONFLICT (phone) DO NOTHING RETURNING id, name, phone, password, address, active, created;
INSERT INTO customers (name, phone, password, address) VALUES ('Fedya', '2', '12345678', '123 Main St') ON CONFLICT (phone) DO NOTHING RETURNING id, name, phone, password, address, active, created;
INSERT INTO customers (name, phone, password, address) VALUES ('Fedya', '3', '12345678', '123 Main St') ON CONFLICT (phone) DO NOTHING RETURNING id, name, phone, password, address, active, created;
INSERT INTO customers (name, phone, password, address) VALUES ('Fedya', '4', '12345678', '123 Main St') ON CONFLICT (phone) DO NOTHING RETURNING id, name, phone, password, address, active, created;
INSERT INTO customers (name, phone, password, address) VALUES ('Fedya', '5', '12345678', '123 Main St') ON CONFLICT (phone) DO NOTHING RETURNING id, name, phone, password, address, active, created;
INSERT INTO customers (name, phone, password, address) VALUES ('Fedya', '6', '12345678', '123 Main St') ON CONFLICT (phone) DO NOTHING RETURNING id, name, phone, password, address, active, created;
SELECT id, name, phone, password, address, active, created FROM customers WHERE name = 'Fedya';

INSERT INTO medicines (name, manafacturer, description, price, pharmacy_name)
VALUES ('med1', 'man1', 'desc1', 1, 'pharm2') RETURNING id, active, created;
INSERT INTO medicines (name, manafacturer, description, price, pharmacy_name)
VALUES ('med1', 'man2', 'desc1', 1, 'pharm2') RETURNING id, active, created;
INSERT INTO medicines (name, manafacturer, description, price, pharmacy_name)
VALUES ('med2', 'man1', 'desc1', 1, 'pharm3') RETURNING id, active, created;
INSERT INTO medicines (name, manafacturer, description, price, pharmacy_name)
VALUES ('med2', 'man1', 'desc1', 1, 'pharm1') RETURNING id, active, created;
INSERT INTO medicines (name, manafacturer, description, price, pharmacy_name)
VALUES ('med3', 'man1', 'desc1', 1, 'pharm3') RETURNING id, active, created;
INSERT INTO medicines (name, manafacturer, description, price, pharmacy_name)
VALUES ('med3', 'man2', 'desc1', 1, 'pharm1') RETURNING id, active, created;