import React from 'react';
import './App.css';
import Navbar from './components/navbar/Navbar';
// import Book from './components/book/book';
import axios from "axios"
import { Route, Routes } from 'react-router-dom';
import Catalogue from './components/pages/catalogue/catalogue';
import Category from './components/pages/category/Category';
import SpecialOffers from "./components/pages/Special-offers/SpecialOffers"
import Bestsellers from './components/pages/Bestsellers/Bestsellers';
import AboutMe from './components/pages/about-me/about-me';
import LoginPage from "./components/pages/login/login";
import Register from "./components/pages/register/register";

// const api = axios.create({
//     baseURL: 'localhost:4001/api/'
// })

class App extends React.Component {
  constructor(props) {
    super(props)
        const userData = localStorage.getItem("userData");
        if (userData !== null) {
            // localStorage.setItem("")
            let greet = document.getElementById("h-greet")
            greet.innerHTML = "Hi, " + userData.Username
        }


  }
s
    componentDidMount() {
        // check if value of id is in the session
      let userData = localStorage.getItem("userData");
      if (userData !== null) {
          let username = userData.Username;
          alert("Your token has not expired");
      } else {
          alert("Your token has expired");
      }
    }

    render () {
    return (

      <div className="root">
          <Navbar></Navbar>
          <div className="App-routes">
              <title>My E-commerce app!</title>
              <Routes>
                <Route path="/home" element={<></>} />
                <Route path="/about-me" element={<AboutMe />} />
                <Route path="/categories" component={<Category />} />
                <Route path="/catalogue" components={<Catalogue />} />
                <Route path="/login" components={<LoginPage />} />
                <Route path="/register" components={<Register />} />
              </Routes>
          </div>
      </div>
    );
  }
}
//  export default NavBar --> import NavBar from "./NavBar"
// jwt-decode
export default App;
