CREATE TABLE BOS_Counter
(
    szCounterId NVARCHAR(50) NOT NULL,
    iLastNumber BIGINT NOT NULL,
    CONSTRAINT PK_BOS_Counter PRIMARY KEY CLUSTERED
        (
         szCounterId ASC
            )
);

CREATE TABLE BOS_Balance
(
    szAccountId NVARCHAR(50) NOT NULL,
    szCurrencyId NVARCHAR(50) NOT NULL,
    decAmount DECIMAL (30, 8) NOT NULL,
    CONSTRAINT PK_BOS_Balance PRIMARY KEY CLUSTERED
        (
         szAccountId ASC, szCurrencyId ASC
            )
);

CREATE TABLE BOS_History
(
    szTransactionId NVARCHAR(50) NOT NULL,
    szAccountId NVARCHAR(50) NOT NULL,
    szCurrencyId NVARCHAR(50) NOT NULL,
    dtmTransaction DATETIME NOT NULL,
    decAmount DECIMAL (30, 8) NOT NULL,
    szNote NVARCHAR(255) NOT NULL,
    CONSTRAINT PK_BOS_History PRIMARY KEY CLUSTERED
        (
         szTransactionId ASC, szAccountId ASC, szCurrencyId ASC
            )
);

-- ------------------
-- INITIALIZATION --
-- ------------------

INSERT INTO BOS_Counter VALUES
                            ('001-COU', 4);

INSERT INTO BOS_Balance VALUES
                            ('000108757484', 'IDR', 34500000.00),
                            ('000108757484', 'USD', 125.8750),
                            ('000109999999', 'IDR', 1250.00),
                            ('000109999999', 'SGD', 2.25),
                            ('000108888888', 'SGD', 125.75);

INSERT INTO BOS_History VALUES
                            ('20201231-00000.00001', '000108757484', 'IDR', GETDATE(), 34500000.00, 'SETOR'),
                            ('20201231-00000.00001', '000108757484', 'SGD', GETDATE(), 125.8750, 'SETOR'),
                            ('20201231-00000.00002', '000109999999', 'IDR', GETDATE(), 1250.00, 'SETOR'),
                            ('20201231-00000.00003', '000109999999', 'SGD', GETDATE(), 128.00, 'SETOR'),
                            ('20201231-00000.00004', '000109999999', 'SGD', GETDATE(), -125.75, 'TRANSFER'),
                            ('20201231-00000.00004', '000108888888', 'SGD', GETDATE(), 125.75, 'TRANSFER');