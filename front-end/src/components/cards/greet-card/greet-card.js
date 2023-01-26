import React from "react";
import "./book.css"

function GreetCard({title}) {
  
  return (
      <>
        <div className="book-cont">
          <div className="content">
            <h2 className="book-title">{title}</h2>
            <div className="book-desc">
                Good day! My name is Oleksandr Matviienko, and I am glad to see you here!
                A bit about myself:
                I am 17 years old, currently open for new oportunities.
                I am ambitious, self-motivated, decisive, organized and hard-working. Believe you'll give me a chance to prove my words.
                I am wishing to work in Start ups! In my opinion, it is best environment for meeting new people, and to test your skills!
                My tech stack:
                I am experienced in JavaScript, Python, Golang, React and I am able to create Web applications, using different
                frameworks and libraries.  Additonally, I have been learning Solidity programming language, Ethereum blockchain concepts, and Web3 for nearly two years.
                I am  holder of good experience of coding applications with Web3 developer tools, such as: Hardhat, Truffle, Ganache, Web3.js, Ethers.js,
                Alchemy, Infura, and Remix ide. At this stage, money is not a priority for me, so I am able to work without any charge, for 40+ hours a week.
            </div>
          </div>
        </div>
        <div className="book-cont">
            <div className="content">
                <h2 className="book-title">{title}</h2>
                <div className="book-desc">
                    Here you can view my resume :
                    <a className="resume" href="https://drive.google.com/file/d/1kD88Stv37ge6KpKT18AUbATqo3rM7zqX/view?usp=sharing">Resume</a>
                </div>
            </div>
        </div>

    </>
  )
}

export default GreetCard;