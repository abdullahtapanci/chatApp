import Navbar from "./Navbar"
import NavbarAfterLoggedIn from "./NavbarAfterLoggedIn";
import { useState, useEffect } from "react";

const Home = () => {

    const [userData, setUserData]=useState("");
    const [isLoggedIn, setIsLoggedIn]=useState(false);
    const [userImageName, setUserImageName] = useState("");
    const [userImageUrl, setUserImageUrl] = useState("");
    const userImageDirectoryPath = "http://localhost:8000/userImages/";

    useEffect( () => {
        const fetchUserData = async () => {
            try{
                const token = localStorage.getItem("jwtToken");
                if(!token){
                    setIsLoggedIn(false)
                    return
                }

                const response = await fetch('http://localhost:8000/api/home', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    }
                });

                if (response.ok) {
                    // Token is valid, extract user data from response
                    const data = await response.json();
                    setUserData(data);
                    if(data.profileImageName === "null"){
                        setUserImageName("profileImageIcon.png")
                        setUserImageUrl(userImageDirectoryPath + "profileImageIcon.png");
                    }else{
                        setUserImageName(data.profileImageName)
                        setUserImageUrl(userImageDirectoryPath + data.profileImageName);
                    }
                    setIsLoggedIn(true);
                } else {
                    console.error('Token validation failed');
                }

            }
            catch (error) {
                console.error('Error validating token:', error);
            }
        }

        fetchUserData();
    },[])

  
    return (
        <div>
            { isLoggedIn ? <NavbarAfterLoggedIn/> : <Navbar/> }
            <h1 className='text-4xl text-center mt-32'>Welcome To Exclusive Chat</h1>
            {isLoggedIn && userData &&
                <div>
                    <div className="border border-gray-500 w-20 h-20 mt-8 rounded-full mx-auto">
                        <img src={userImageUrl} alt={userImageName} className="w-full h-full object-cover rounded-full"></img>
                    </div>
                    <div className="flex justify-center mt-8">
                        <h3 className='text-xl text-center p-2'>{userData.firstName}</h3>
                        <h3 className='text-xl text-center p-2'>{userData.lastName}</h3>
                    </div>
                    <h3 className='text-xl text-center p-2'>{userData.userName}</h3>
                </div>
                
            }
        </div>
    );
}

export default Home;
