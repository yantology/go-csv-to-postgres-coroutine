CREATE TABLE IF NOT EXISTS domain (
    GlobalRank INTEGER,
    TldRank INTEGER,
    Domain VARCHAR(255),
    TLD VARCHAR(255),
    RefSubNets INTEGER,
    RefIPs INTEGER,
    IDN_Domain VARCHAR(255),
    IDN_TLD VARCHAR(255),
    PrevGlobalRank INTEGER,
    PrevTldRank INTEGER,
    PrevRefSubNets INTEGER,
    PrevRefIPs INTEGER
);