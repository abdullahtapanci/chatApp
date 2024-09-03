import NavbarAfterLoggedIn from "./NavbarAfterLoggedIn"
import { useState, useEffect } from "react";

const FindAFriend = () => {

    const [searchQuery, setSearchQuery] = useState("");
    const [foundUsers, setFoundUsers] = useState();

    const userImageDirectoryPath = "http://localhost:8000/userImages/";

    const handleSearchInputChange = (e) => {
        setSearchQuery(e.target.value)
    }

    //find users according to search query
    useEffect(()=>{
        const fetchUserData = async () => {
            try{
                const token = localStorage.getItem("jwtToken");

                const response = await fetch('http://localhost:8000/api/findUser', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify(searchQuery)
                });

                if (response.ok) {
                    // Token is valid, extract user data from response
                    const data = await response.json();
                    const foundUsersData = JSON.parse(data.foundUsers);
                    setFoundUsers(foundUsersData);
                } else {
                    console.error('Couldnt fetch user data');
                }
            }
            catch (error) {
                console.error('Error fetching user data :', error);
            }
        }

        fetchUserData();
    },[searchQuery])  

    const createFriendship = async (id) => {
        try{
            const token = localStorage.getItem("jwtToken");
            const requestData = {
                friendId: id
            };

            const response = await fetch('http://localhost:8000/api/createFriendship', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(requestData)
            });

            if (response.ok) {
                window.location.href = '/chat';
                
            } else {
                console.error('Couldnt be friend');
            }
        }
        catch (error) {
            console.error('Error being friend :', error);
        }
    }

    return (
        <div>
            <NavbarAfterLoggedIn></NavbarAfterLoggedIn>
            <h1 className="text-center mt-8 text-2xl">Find Your Friend And Say Hi!</h1>

            <div className="mt-8 flex justify-center items-center">
                <input
                    type="text"
                    placeholder="Search..."
                    value={searchQuery}
                    onChange={handleSearchInputChange}
                    className="w-1/2 p-2 border rounded-lg focus:outline-none focus:ring focus:border-blue-500"
                />
            </div>
            <div className="mt-8 p-2">
                {
                    foundUsers ? ( 
                        foundUsers.map(user => (
                            <div key={user.userName} className="flex justify-between">
                                <div className="flex items-center m-4">
                                    <div className="border border-gray-500 w-12 h-12 rounded-full mr-4">
                                    {user.imageName !== "null" &&(
                                        <img src={userImageDirectoryPath + user.imageName} alt={user.imageName} className="w-full h-full object-cover rounded-full"></img>
                                    )}
                                    {user.imageName === "null" &&(
                                        <img src={userImageDirectoryPath + "profileImageIcon.png"} alt={"profileImageIcon"} className="w-full h-full object-cover rounded-full"></img>
                                    )}
                                    </div>
                                    <div>
                                        <p className="font-bold" >{user.firstName} {user.lastName}</p>
                                        <p>{user.userName}</p>
                                    </div>
                                </div>
                                <button className="px-4 py-2 m-4 border border-gray-500 rounded hover:bg-blue-500 hover:text-white" onClick={() => createFriendship(user.id)}>Chat</button>
                            </div>
                        ))
                    ) : (
                        <div></div>
                    )
                }
            </div>

        </div>
    )
}

export default FindAFriend;