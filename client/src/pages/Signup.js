import Navbar from "./Navbar"
import { useState } from "react";

const Signup = () =>{
    //create userData State
    const [userData,setUserData]=useState({
        firstName:"",
        lastName:"",
        userName:"",
        email:"",
        password:"",
    });
    
    const handleChange = (e) => {
        const { name, value } = e.target;
        setUserData((prevData) => ({ ...prevData, [name]: value }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        // send data to server to check user if it is in database
        try{
            const response = await fetch('http://localhost:8000/api/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });
    
            // Check if the request was successful
            if (response.ok) {
                const responseData = await response.json();
                console.log('Server response:', responseData);
                window.location.href = '/login';
                alert('Signup successful!');
            } else {
                console.error('Failed to submit the form');
            }
        }
        catch{
            
        }

    };

    return(
        <div>
            <Navbar></Navbar>
            <h1 className="text-center text-3xl mt-10">Sign Up</h1>
            <form onSubmit={handleSubmit} className="flex flex-col">
                <div className="flex flex-col">
                    <input 
                        type="text" id="firstName" name="firstName" value={userData.firstName} placeholder="Name" onChange={handleChange} required
                        className="w-1/4 m-2 mt-10 mx-auto p-2 border rounded"
                    >
                    </input>
                    <input 
                        type="text" id="lastName" name="lastName" value={userData.lastName} placeholder="Surname" onChange={handleChange} required
                        className="w-1/4 m-2 mx-auto p-2 border rounded"
                    >
                    </input>
                    <input 
                        type="text" id="userName" name="userName" value={userData.userName} placeholder="Username" onChange={handleChange} required
                        className="w-1/4 m-2 mx-auto p-2 border rounded"
                    >
                    </input>
                    <input 
                        type="email" id="email" name="email" value={userData.email} placeholder="Email" onChange={handleChange} required
                        className="w-1/4 m-2 mx-auto p-2 border rounded"
                    >
                    </input>
                    <input 
                        type="password" id="password" name="password" value={userData.password} placeholder="Password" onChange={handleChange} required
                        className="w-1/4 m-2 mb-10 mx-auto p-2 border rounded"
                    >
                    </input>
                </div>
                <button type="submit" className=" max-w-xl mx-auto p-2 border rounded hover:bg-blue-500 hover:text-white">Submit</button>
            </form>
        </div>
    )
}

export default Signup;