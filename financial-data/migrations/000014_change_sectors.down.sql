DELETE * from sectors;

INSERT INTO sectors (id, name) VALUES
    (1, 'Oils'),
    (2, 'Finance'),
    (3, 'Technology'),
    (4, 'Telecom'),
    (5, 'Metals'),
    (6, 'Mining'),
    (7, 'Utilities'),
    (8, 'RealEstate'),
    (9, 'ConsumerStaples'),
    (10, 'ConsumerDiscretionary'),
    (11, 'Healthcare'),
    (12, 'Industrial'),
    (13, 'Energy'),
    (14, 'Materials'),
    (15, 'Transportation'),
    (16, 'Agriculture'),
    (17, 'Chemicals'),
    (18, 'Construction'),
    (19, 'Retail')
ON CONFLICT (id) DO NOTHING;