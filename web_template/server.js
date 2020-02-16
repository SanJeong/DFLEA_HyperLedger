// ExpressJS Setup
const express = require('express');
const app = express();
var bodyParser = require('body-parser');

// Hyperledger Bridge
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'first-network' ,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Constants
const PORT = 8800;
const HOST = '0.0.0.0';

// view engine setup : ejs
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'ejs');

// use static file
app.use(express.static(path.join(__dirname, 'public')));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// DB
var mongoose = require('mongoose');
// indexx 부여를 위한
var autoIncrement = require('mongoose-auto-increment');
mongoose.set('useCreateIndex', true);
mongoose.connect('mongodb://localhost:27017/web_template', { useNewUrlParser: true,  useUnifiedTopology: true  })
var db = mongoose.connection
autoIncrement.initialize(db);

db.on('error', function(){
    console.log('MongoDB connection failed!')
})
db.once('open', function(){
    console.log('MongoDB connection success!')
})


// Passport 로그인 및 회원가입 로직 
var passport = require('passport')
var session = require('express-session')
// flash는 session이 필요하므로 반드시 session 아래에 정의해야 함
var flash = require('connect-flash');
app.use(session({
  secret: 'keyboard cat',
  resave: false,
  saveUninitialized: true
}))
app.use(passport.initialize())
app.use(passport.session())
app.use(flash())

//파일 업로드
var multer = require('multer');
const upload = multer({
  storage: multer.diskStorage({
    destination: function (req, file, cb) {
      cb(null, 'upload/');
    },
    filename: function (req, file, cb) {
      cb(null, new Date().valueOf() + path.extname(file.originalname));
    }
  }),
});

// Router
var indexRouter = require('./router/index');
var loginRouter = require('./router/login');
var joinRouter = require('./router/join');
var logoutRouter = require('./router/logout');
var itemsRouter = require('./router/items');
var newRouter = require('./router/new');
var rateRouter = require('./router/rate');
var mypageRouter = require('./router/mypage');

app.use('/', indexRouter);
app.use('/login', loginRouter);
app.use('/join', joinRouter);
app.use('/logout', logoutRouter);
app.use('/items', itemsRouter);
app.use('/new', newRouter);
app.use('/rate', rateRouter);
app.use('/mypage', mypageRouter);


// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
