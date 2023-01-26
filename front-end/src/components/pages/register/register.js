import React, {useState} from "react";
import Navbar from "../../navbar/Navbar";
import "./register.css";
import axios from "axios";

export default function Register() {
    const [firstname, setFirstname] = useState("")
    const [lastname, setLastname] = useState("")
    const [email, setEmail] = useState("")

    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [repPassword, setRepPassword] = useState("")

    const regMessages = document.getElementById("")

    function CheckPassword() {
        if (password !== repPassword) {
            regMessages.classList.remove("d-none")
            regMessages.innerHTML = "Password and repeated password should be equal."
            return false
        }
        return true
    }

    function ShowErrorMessges(msg) {
        regMessages.classList.add("alert-danger")
        regMessages.innerHTML = msg
    }

    function PostRegister() {
        if (!CheckPassword()) {
            return
        }
        const config = {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            }
        }
        const data = {
            "firstname": firstname,
            "lastname": lastname,
            "email": email,
            "password": password
        }

        axios.post("localhost:4001/api/signup", data,  config)
            .then(res => res.json)
            .then(data => {
                console.log(data);
                if (data.error === false) {
                    ShowErrorMessges(data.message);
                } else {
                    window.location.href = "/login";
                }
        })
    }

    return (
        <>
        <Navbar></Navbar>
            <div className="login-box">
                <h2>Login</h2>
                <form method="post" action="">
                    <div className="user-box">
                        <input value={firstname} onChange={e => setFirstname(e.target.message)} type="text" name="firstname" required="" />
                        <label>Firstname</label>
                    </div>
                    <div className="user-box">
                        <input value={lastname} onChange={e => setLastname(e.target.message)} type="text" name="lastname" required="" />
                        <label>Lastname</label>
                    </div>
                    <div className="user-box">
                        <input value={username} onChange={e => setUsername(e.target.message)} type="text" name="username" required="" />
                        <label>Username</label>
                    </div>
                    <div className="user-box">
                        <input value={email} onChange={e => setEmail(e.target.message)} type="text" name="email" required="" />
                        <label>Email</label>
                    </div>
                    <div className="user-box">
                        <input value={password} onChange={e => setPassword(e.target.message)} type="password" name="password" required="" />
                        <label>Password</label>
                    </div>
                    <div className="user-box">
                        <input value={repPassword} onChange={e => setRepPassword(e.target.message)} type="password" name="password" required="" />
                        <label>Repeat Password</label>
                    </div>
                    <a href="javascript:void(0)" onClick="PostRegister()">
                        <span></span>
                        <span></span>
                        <span></span>
                        <span></span>
                        Submit
                    </a>
                </form>
            </div>
        </>
    )
}