INSERT INTO customers (name, phone, password, address) VALUES 
    ('Firya', '1', '12345678', '123 Main St'),
    ('Misha', '2', '12345678', '123 Main St'),
    ('Vasya', '3', '12345678', '123 Main St'),
    ('Kirya', '4', '12345678', '123 Main St')
ON CONFLICT (phone) DO NOTHING
RETURNING id, name, phone, password, address, active, created;
INSERT INTO customers (name, phone, password, address, is_admin) VALUES 
    ('mr. Fedya', '5', '12345678', '123 Main St', TRUE)
ON CONFLICT (phone) DO NOTHING
RETURNING id, name, phone, password, address, active, created;
SELECT id, name, phone, password, address, active, created FROM customers WHERE name = 'Fedya';

INSERT INTO medicines (name, manafacturer, description, price, pharmacy_name) VALUES
    ('med1', 'man1', 'desc1', 1, 'pharm1'),
    ('med2', 'man1', 'desc2', 2, 'pharm2'),
    ('med3', 'man1', 'desc3', 3, 'pharm3'),
    ('med1', 'man2', 'desc4', 4, 'pharm2'),
    ('med2', 'man2', 'desc5', 5, 'pharm2'),
    ('med3', 'man2', 'desc6', 6, 'pharm1'),
    ('med1', 'man3', 'desc7', 7, 'pharm3'),
    ('med2', 'man3', 'desc8', 8, 'pharm3'),
    ('med3', 'man3', 'desc9', 9, 'pharm2')
RETURNING id, active, created;