var express = require('express');
var router = express.Router();
var passport = require('passport');
var User = require('../model/user');

passport.use(User.createStrategy());
passport.serializeUser(User.serializeUser());
passport.deserializeUser(User.deserializeUser());

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname,'..', '..', 'first-network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

async function cc_call(fn_name, args){

  const walletPath = path.join(process.cwd(),'wallet');
  const wallet = new FileSystemWallet(walletPath);
  console.log(`Wallet path: ${walletPath}`);

  const userExists = await wallet.exists('admin');
  if (!userExists) {
      console.log('An identity for the user "user1" does not exist in the wallet');
      console.log('Run the registerUser.js application before retrying');
      return;
  }
  const gateway = new Gateway();
  await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: false } });
  const network = await gateway.getNetwork('mychannel');
  const contract = network.getContract('dflea');

  var result;

  console.log(`saarc:fn:${fn_name}, args:${args}`);

  if(fn_name == 'addUser')
      result = await contract.submitTransaction('addUser', args[0],args[1]);
  else if( fn_name == 'addDataset')
      result = await contract.submitTransaction('addDataset', args[0],args[1],args[2],args[3],args[4]);
  else if( fn_name == 'addPurchase')
      result = await contract.submitTransaction('addPurchase', args[0],args[1],args[2],args[3]);
  else if( fn_name == 'transferPurchase')
      result = await contract.submitTransaction('transferPurchase', args[0],args[1]);
  else if(fn_name == 'readUser')
      result = await contract.evaluateTransaction('readUser', args);
  else if(fn_name == 'readDataset')
      result = await contract.evaluateTransaction('readDataset', args);
  else if(fn_name == 'readDatasetPrivateDetails')
      result = await contract.evaluateTransaction('readDatasetPrivateDetails', args);
  else if(fn_name == 'readPurchase')
      result = await contract.evaluateTransaction('readPurchase', args);
  else
      result = 'not supported function'

  return result;
}

/* GET home page. */
router.get('/', function(req, res, next) {
    res.render('join', {title: "Join"});
});

// 회원가입 로직
router.post('/',  function(req, res, next) {
    console.log(req.body.email);
    console.log(req.body.name);
    console.log(req.body.password);
    const handleRegister = (err, user)=>{
      console.log(err)
    }

    var email = req.body.email;
    var name = req.body.name;

    var args = [email, name];

    console.log(args);
    // 블록체인에 등록
    result = cc_call('addUser', args)
    const myobj = {result: "success"}

    // DB에 회원등록
    User.register(new User({name: req.body.name, email: req.body.email}), req.body.password, function(err) {
      if (err) {
        console.log('error while user register!', err);
        return next(err);
      }
      console.log('회원가입 성공');
      res.redirect('/');
    });
  })

module.exports = router;
