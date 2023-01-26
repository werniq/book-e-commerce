import React from "react";
import "./bigger.css"

function Bigger({description}) {
  
  return (
        <div className="big-cont">
          <div className="content">
            <h2 className="big-title">About me:</h2>
            <div className="big-desc">
                {description}
            </div>
          </div>
        </div>
  )
}

export default Bigger;