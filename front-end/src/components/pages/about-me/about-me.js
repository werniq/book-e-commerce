import React from "react";
import Bigger from "../../cards/greet-card/bigger-greet";
import GreetCard from "../../cards/greet-card/greet-card";
import "./index.css";
import Navbar from "../../navbar/Navbar";

function AboutMe() {
	return (
    <div className="home">
		<GreetCard title="A bit about me..">
		</GreetCard>

		{/*<GreetCard title="A bit about me..">*/}
		{/*</GreetCard>*/}
		<div className="space-1">

		</div>
		</div>
    )
}

export default AboutMe;