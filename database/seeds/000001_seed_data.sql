-- 1. USERS
INSERT INTO users (id, full_name, email, password_hash, phone_number, role, is_active)
VALUES 
    ('a1111111-1111-1111-1111-111111111111', 'Super Admin MKP', 'admin@mkp.com', '$2a$10$.gEylANSML1gSh7/QkCaGe6niOpRD0NETl4Y5RThXS1LnbUJr4Czq', '081200000001', 'ADMIN', true),
    ('a2222222-2222-2222-2222-222222222222', 'Cinema Staff Jakarta', 'staff@mkp.com', '$2a$10$.gEylANSML1gSh7/QkCaGe6niOpRD0NETl4Y5RThXS1LnbUJr4Czq', '081200000002', 'STAFF', true),
    ('a3333333-3333-3333-3333-333333333333', 'Budi Santoso (Customer)', 'customer@mkp.com', '$2a$10$.gEylANSML1gSh7/QkCaGe6niOpRD0NETl4Y5RThXS1LnbUJr4Czq', '081200000003', 'CUSTOMER', true),
    ('a4444444-4444-4444-4444-444444444444', 'Siti Rahma (Customer)', 'siti.rahma@mkp.com', '$2a$10$.gEylANSML1gSh7/QkCaGe6niOpRD0NETl4Y5RThXS1LnbUJr4Czq', '081200000004', 'CUSTOMER', true)
ON CONFLICT (email) DO NOTHING;

-- 2. CITIES
INSERT INTO cities (id, name, province)
VALUES 
    ('b1111111-1111-1111-1111-111111111111', 'Jakarta Pusat', 'DKI Jakarta'),
    ('b2222222-2222-2222-2222-222222222222', 'Surabaya', 'Jawa Timur'),
    ('b3333333-3333-3333-3333-333333333333', 'Denpasar', 'Bali'),
    ('b4444444-4444-4444-4444-444444444444', 'Bandung', 'Jawa Barat'),
    ('b5555555-5555-5555-5555-555555555555', 'Semarang', 'Jawa Tengah')
ON CONFLICT (name) DO NOTHING;

-- 3. CINEMAS (CABANG NASIONAL)
INSERT INTO cinemas (id, city_id, name, address, phone_number, is_active)
VALUES 
    ('c1111111-1111-1111-1111-111111111111', 'b1111111-1111-1111-1111-111111111111', 'MKP Grand Indonesia XXI', 'Jl. M.H. Thamrin No.1, Jakarta Pusat', '021-23580001', true),
    ('c2222222-2222-2222-2222-222222222222', 'b2222222-2222-2222-2222-222222222222', 'MKP Galaxy Mall Cinema', 'Jl. Dharmahusada Indah Timur No.35, Surabaya', '031-5937100', true),
    ('c3333333-3333-3333-3333-333333333333', 'b3333333-3333-3333-3333-333333333333', 'MKP Beachwalk Bali Premiere', 'Jl. Pantai Kuta, Denpasar, Bali', '0361-8464888', true)
ON CONFLICT DO NOTHING;

-- 4. STUDIOS
INSERT INTO studios (id, cinema_id, name, studio_type, total_rows, total_cols, total_seats, is_active)
VALUES 
    ('d1111111-1111-1111-1111-111111111111', 'c1111111-1111-1111-1111-111111111111', 'Studio 1 Regular', 'REGULAR', 5, 8, 40, true),
    ('d2222222-2222-2222-2222-222222222222', 'c1111111-1111-1111-1111-111111111111', 'Studio 2 IMAX Laser', 'IMAX', 6, 10, 60, true),
    ('d3333333-3333-3333-3333-333333333333', 'c1111111-1111-1111-1111-111111111111', 'Studio 3 Premiere Lounge', 'PREMIERE', 4, 6, 24, true),
    ('d4444444-4444-4444-4444-444444444444', 'c2222222-2222-2222-2222-222222222222', 'Studio 1 Regular', 'REGULAR', 5, 8, 40, true)
ON CONFLICT DO NOTHING;

-- 5. SEATS (Contoh Baris A-C, Kolom 1-4 untuk Studio 1 Regular)
INSERT INTO seats (studio_id, seat_number, row_name, col_number, is_active)
SELECT 
    'd1111111-1111-1111-1111-111111111111'::UUID,
    r.row || c.col::TEXT,
    r.row,
    c.col,
    true
FROM 
    (VALUES ('A'), ('B'), ('C'), ('D'), ('E')) AS r(row)
    CROSS JOIN (VALUES (1), (2), (3), (4), (5), (6), (7), (8)) AS c(col)
ON CONFLICT DO NOTHING;

-- 6. MOVIES
INSERT INTO movies (id, title, description, duration_minutes, genre, age_rating, poster_url, release_date, is_active)
VALUES 
    (
        'e1111111-1111-1111-1111-111111111111',
        'Dune: Part Two',
        'Paul Atreides bersatu dengan Chani dan kaum Fremen untuk membalas dendam terhadap para konspirator yang menghancurkan keluarganya.',
        166,
        'Sci-Fi / Adventure',
        '13+',
        'https://image.tmdb.org/t/p/w500/dune2.jpg',
        '2024-03-01',
        true
    ),
    (
        'e2222222-2222-2222-2222-222222222222',
        'Agak Laen 2',
        'Empat sekawan penjaga wahana rumah hantu kembali menghadapi petualangan misteri konyol yang jauh lebih mengocok perut dan menegangkan.',
        119,
        'Comedy / Horror',
        '13+',
        'https://image.tmdb.org/t/p/w500/agaklaen2.jpg',
        '2025-01-15',
        true
    ),
    (
        'e3333333-3333-3333-3333-333333333333',
        'Avatar: Fire and Ash',
        'Petualangan epik Jake Sully dan Neytiri berlanjut di wilayah suku abu Pandora yang keras dan penuh konflik baru.',
        190,
        'Action / Sci-Fi / Fantasy',
        'SU',
        'https://image.tmdb.org/t/p/w500/avatar3.jpg',
        '2025-12-19',
        true
    )
ON CONFLICT DO NOTHING;

-- 7. SCHEDULES (JADWAL TAYANG)
INSERT INTO schedules (id, movie_id, studio_id, start_time, end_time, price, status)
VALUES 
    (
        'f1111111-1111-1111-1111-111111111111',
        'e1111111-1111-1111-1111-111111111111',
        'd1111111-1111-1111-1111-111111111111',
        CURRENT_DATE + TIME '13:00:00',
        CURRENT_DATE + TIME '15:50:00',
        50000.00,
        'SCHEDULED'
    ),
    (
        'f2222222-2222-2222-2222-222222222222',
        'e1111111-1111-1111-1111-111111111111',
        'd1111111-1111-1111-1111-111111111111',
        CURRENT_DATE + TIME '16:30:00',
        CURRENT_DATE + TIME '19:20:00',
        50000.00,
        'SCHEDULED'
    ),
    (
        'f3333333-3333-3333-3333-333333333333',
        'e2222222-2222-2222-2222-222222222222',
        'd1111111-1111-1111-1111-111111111111',
        CURRENT_DATE + TIME '20:00:00',
        CURRENT_DATE + TIME '22:15:00',
        60000.00,
        'SCHEDULED'
    ),
    (
        'f4444444-4444-4444-4444-444444444444',
        'e1111111-1111-1111-1111-111111111111',
        'd2222222-2222-2222-2222-222222222222',
        CURRENT_DATE + TIME '14:00:00',
        CURRENT_DATE + TIME '16:50:00',
        95000.00,
        'SCHEDULED'
    )
ON CONFLICT DO NOTHING;

-- 8. BOOKINGS & PAYMENTS (Contoh Transaksi Selesai)
INSERT INTO bookings (id, booking_code, user_id, schedule_id, total_tickets, total_amount, status, expires_at)
VALUES 
    (
        '11111111-bbbb-bbbb-bbbb-111111111111',
        'MKP-BK-20250915-001',
        'a3333333-3333-3333-3333-333333333333',
        'f1111111-1111-1111-1111-111111111111',
        2,
        100000.00,
        'PAID',
        CURRENT_TIMESTAMP + INTERVAL '1 hour'
    )
ON CONFLICT DO NOTHING;

INSERT INTO payments (id, booking_id, transaction_ref, payment_method, amount, status, paid_at, expires_at, payment_gateway_response)
VALUES 
    (
        '11111111-cccc-cccc-cccc-111111111111',
        '11111111-bbbb-bbbb-bbbb-111111111111',
        'TRX-PGW-20250915-99881',
        'QRIS',
        100000.00,
        'SUCCESS',
        CURRENT_TIMESTAMP - INTERVAL '30 minutes',
        CURRENT_TIMESTAMP + INTERVAL '10 minutes',
        '{"gateway": "Xendit", "channel": "QRIS", "status": "COMPLETED", "currency": "IDR"}'::JSONB
    )
ON CONFLICT DO NOTHING;
