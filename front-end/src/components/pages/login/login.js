import React, {useState} from "react";
import Navbar from "../../navbar/Navbar";
import axios from "axios";


export default function LoginPage() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")

    function ShowErrorMsg(e) {
        const card = document.getElementById("login-messages")
        card.classList.remove("d-none")
    }

    function Login() {
        const config = {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            }
        }
        const data = {
            "email": email,
            "password": password
        }
        axios.post("localhost:4001/api/authenticate", data, config)
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    const msgs = document.getElementById("login-messages");
                    msgs.classList.remove('d-none')
                    msgs.innerHTML = data.message
                    return
                } else {
                    let userData = {
                        "userID":data.UserId,
                        "expriy": data.TokenExpiry,
                        "role": data.Role,
                        "token": data.Token,
                        "username": data.Username
                    }
                    localStorage.setItem("userData", JSON.stringify(userData))
                }
            })
         }

    return (
        <>
            <head>
                <link rel="stylesheet” href=”https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css”rel=”nofollow” integrity=”sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm” crossorigin=”anonymous" />
            </head>
            <Navbar></Navbar>
            <div className="login-box">
                <h2>Login</h2>
                { localStorage.getItem("userID") !== null
                    ?
                <form>
                    <div className="alert alert-danger text-center d-none" id="login-messages"></div>

                    <div className="user-box">
                        <input type="text" value={email} onChange={e => setEmail(e.target.value)} name="email" required="" />
                            <label>Username</label>
                    </div>
                    <div className="user-box">
                        <input type="password" value={password} onChange={e => setPassword(e.target.value)}  name="password" required="" />
                            <label>Password</label>
                    </div>
                    <a href="javascript:void(0)" onClick="Login()">
                        <span></span>
                        <span></span>
                        <span></span>
                        <span></span>
                        Submit
                    </a>
                </form>
                    : <div className="loginned">You are already logged in </div>
                }
            </div>
        </>
    )
}