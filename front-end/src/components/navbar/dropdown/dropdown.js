import React from "react";
import "./dropdown.css";

export default function Dropdown({...elements}) {
    return (
        <div class="dropdown-content">
        			{elements.Map(element => (
        			  <div>
        			      <a>{element}</a>
        			  </div>
        			))}
        </div>
    
    )
}