import React, { useEffect, useState } from "react";
import { NavbarNav } from "./NavbarElements";
import "./style.css";
import Web3 from "web3"
import "./dropdown/dropdown.css"
import {Link, NavLink} from "react-router-dom";

export default function Navbar() {
		async function EthConnect() {
			const button = document.getElementById("connect")
			if (window.ethereum) {
					try {
							const account = await window.ethereum.request({method: 'eth_requestAccounts'})
							window.web3 = new Web3(window.ethereum)
							console.log(account)
							button.innerHTML = "Connected"
							let addr = document.getElementById("eth-address")
							addr.innerHTML = account.address
							button.style.display = "none"
					} catch (e) {
							console.log(e)
					}
				}
		}
	let userData = localStorage.getItem("userData")
	let userId = userData.userID;

	function Logout() {
		localStorage.removeItem("token")
		localStorage.removeItem("token_expiry")
		window.location.href = "/login"
	}

	return (
      	<nav>
			<ul>
      	  		<li> <NavLink className="nav-link" to="/home">Home page</NavLink></li>
				<li> <NavLink className="nav-link" to="/special-offers">Special offers</NavLink></li>
      	  		<li> <NavLink className="nav-link" to="/catalogue">Catalogue</NavLink></li>
				<li> <NavLink className="nav-link" to="/best-sellers">Best sellers</NavLink></li>
				<li> <NavLink className="nav-link" to="/categories" >Categories</NavLink></li>
				<li> <button className="metamask-but" id="connect" onClick={EthConnect}>Connect</button></li>
				<li> <NavLink to="/profile" className="nav-profile-link">Profile</NavLink></li>
				{ 	userData === null ?
					<li> <NavLink to="/profile" className="nav-login-link">Register</NavLink></li>
						:
					<li> <NavLink to="/login" className="nav-login-link">Login</NavLink></li>
				}
			</ul>
		</nav>
    )
}