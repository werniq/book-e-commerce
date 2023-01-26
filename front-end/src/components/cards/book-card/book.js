import React from "react";
import "./book.css"

function CardBook({title, description, author, price, image, href}) {
  let Foto = `../Images/${image}`
  
  return (
        <div className="book-cont">
        
          <div className="content">
            <h2 className="book-title">{title}</h2>
            <div className="book-desc">{description}</div>
            <div className="card-footer">
              <a href="" className="read-mr" onClick="redirect">Order</a>
              <a href="" className="read-mr1" onClick="redirect">Read more</a>
              <h2 className="card-price">{price} </h2>
             </div>
            </div>
            <a href={href}><img width="275" height="350" src={Foto} alt="image" /></a>
        </div>
  )
}

export default CardBook;