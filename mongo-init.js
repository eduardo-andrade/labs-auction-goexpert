db = db.getSiblingDB('auctions');

db.createUser({
  user: 'admin',
  pwd: 'admin',
  roles: [{
    role: 'readWrite',
    db: 'auctions'
  }]
});

db.createCollection('auctions');
db.auctions.createIndex({ "end_time": 1 });

db.createCollection('bids');
db.bids.createIndex({ "auction_id": 1 });

db.createCollection('users');