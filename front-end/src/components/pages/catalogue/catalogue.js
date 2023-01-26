import axios from "axios";
import React, {useState, useEffect} from "react";
import CardBook from "../../cards/book-card/book";


export default function Catalogue() {
    // const [products, setProducts] = useState(null)
    const [categories, setCategories] = useState(null)
    const [error, setError] = useState(null)


		useEffect(() => {
			axios.get("localhost:4001/api/categories")
			.then((categories) => {
				setCategories(categories.data)
				console.log(categories)
			}).catch(error => {
				setError(error)
			})
		}, [])



		const RenderBooks = (Books) => {
  	    	return (
					<div className="main-section">
						{Books.map((book, index) => (
  	          			<CardBook
							key={index}
							title={book.title}
							author={book.author}
							description={book.description}
							price={book.price}
							image={book.image}
							href={book.href}>
						</CardBook>
						))}
					</div>
			)
		}
		let books;
	useEffect(() => {
		axios.get("localhost:4001/api/get-products")
			.then((response) => {
				books = response.data
				// products =
				console.log(books)
			}).catch (error => {
			setError(error);
		})
	}, []);

	if (error) return `Error: ${error.message}`;
	if (!books) return <video style={{width: '50px', float:'center'}} src="https://icons8.com/icon/Lji74RNucC7l/hourglass"></video>
	return (
		<div className="catelogue">
			{RenderBooks(books)}
		</div>
	)
	/*
 	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Price       int       `json:"amount"`
	Image       string    `json:"image"`
	IsAvaliable bool      `json:"is_avaliable"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	*/
}