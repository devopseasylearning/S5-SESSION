const express = require('express')
const app = express()
const port = 3000
const path = require('path')
const axios = require('axios')
const cookieParser = require("cookie-parser");
const bodyParser = require('body-parser');

app.use('/', express.static(path.join(__dirname, 'public/static')))
app.set('view engine', 'ejs');
app.use(bodyParser.urlencoded({extended: true}));
app.use(cookieParser());

const jwt = require('jsonwebtoken');

function authenticateToken(req, res, next) {
    if (!req.cookies.token) {
        res.redirect("/login");
        return
    }
    jwt.verify(req.cookies.token, "xco0sr0fh4e52x03g9mv", function (err, decoded) {
        if (err) {
            console.log(err)
            res.redirect("/login")
            return
        }
        next()
    })
}

app.get("/health", (req,res) => {
    res.sendStatus(200)
} )

app.get('/login', (req, res) => {
    res.sendFile(path.join((__dirname + "/public//login.html")))
})
app.get('/signup', (req, res) => {
    res.sendFile(path.join((__dirname + "/public//signup.html")))
})
app.post('/login', (loginreq, loginres) => {
    axios
        .post('http://' + process.env.AUTH_HOST + ':' + process.env.AUTH_PORT + '/users/' + loginreq.body.username, {
            user_name: loginreq.body.username,
            user_password: loginreq.body.password
        })
        .then((res) => {
            loginres.writeHead(302, {
                "Set-Cookie": "token=" + res.data.JWT + "; HttpOnly",
                "Access-Control-Allow-Credentials": "true",
                "Location": "/"
            })
                .send().end();
        })
        .catch((error) => {
            loginres.redirect("/login?error=invalidcreds")
            console.error(error)
        })
})
app.post("/signup", (signupreq, signupres) => {
    axios.post('http://' + process.env.AUTH_HOST + ':' + process.env.AUTH_PORT + '/users/', {
        user_name: signupreq.body.username,
        user_password: signupreq.body.password
    }).then((res) => {
        signupres.redirect("/login")
    }).catch((error) => {
        console.log(error)
        signupres.redirect('/signup?error=userexists')
    })
})

app.get('/', authenticateToken, (req, res) => {
    res.sendFile(path.join((__dirname + "/public//index.html")))
})
app.get("/logout",(req,res) => {
    res.writeHead(302, {
        "Set-Cookie": `token=; HttpOnly; path=/; max-age=0`,
        "Location": "/login"
    });
    res.end();
})
app.get("/weather/:city",(req,res) => {
    axios.get('http://' + process.env.WEATHER_HOST + ':' + process.env.WEATHER_PORT + '/' + req.params.city)
        .then((response) => {
            res.json(response.data)
        }).catch((error) => {
            console.log(error)
        })
})
app.listen(port, () => {
    console.log(`Weather app listening at http://localhost:${port}`)
})