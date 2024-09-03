import { useEffect, useState } from "react";
import NavbarAfterLoggedIn from "./NavbarAfterLoggedIn";

const Profile = () => {

    const [userData,setUserData]=useState({
        firstName:"",
        lastName:"",
        userName:"",
        email:"",
    });


    const [userImageName, setUserImageName] = useState("");
    const [userImageUrl, setUserImageUrl] = useState("");
    const userImageDirectoryPath = "http://localhost:8000/userImages/";

    //get user ınformations
    useEffect(()=>{
        const fetchUserData = async () => {
            try{
                const token = localStorage.getItem("jwtToken");

                const response = await fetch('http://localhost:8000/api/profile', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    }
                });

                if (response.ok) {
                    // Token is valid, extract user data from response
                    const data = await response.json();
                    setUserData({
                        firstName: data.firstName,
                        lastName: data.lastName,
                        userName: data.userName,
                        email: data.email
                    });
                    if(data.profileImageName === "null"){
                        setUserImageName("profileImageIcon.png")
                        setUserImageUrl(userImageDirectoryPath + "profileImageIcon.png");
                    }else{
                        setUserImageName(data.profileImageName)
                        setUserImageUrl(userImageDirectoryPath + data.profileImageName);
                    }
                } else {
                    console.error('Couldnt fetch user data');
                }

            }
            catch (error) {
                console.error('Error fetching user data :', error);
            }
        }

        fetchUserData();
    },[])  

    const handleChange = (e) =>{
        const { name, value } = e.target;
        setUserData((prevData) => ({ ...prevData, [name]: value }));
    }

    //update user profile
    const handleUpdate = async (e) =>{
        e.preventDefault();

        // send data to server to check user if it is in database
        try{
            const token = localStorage.getItem("jwtToken");

            const response = await fetch('http://localhost:8000/api/updateProfileInfo', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}` // Include the token in the Authorization header
                },
                body: JSON.stringify(userData)
            });
    
            // Check if the request was successful
            if (response.ok) {
                const data = await response.json();
                    setUserData({
                        firstName: data.firstName,
                        lastName: data.lastName,
                        userName: data.userName,
                        email: data.email
                    });
                localStorage.setItem('token', data.token);
                alert("Your informations updated successfuly")
                window.location.reload();
            } else {
                console.error('Failed to submit the form');
            }
        }
        catch{
            
        }
    }


    const [imageFile, setImageFile] = useState(null);

    const handleImageChange = (e) => {
        setImageFile(e.target.files[0]);
    };

    //edit user image
    const editUserPhoto = async (e) => {
        e.preventDefault();
        const formData = new FormData();
        formData.append('image', imageFile);

        console.log(formData)

        try {
            const token = localStorage.getItem("jwtToken");
            const response = await fetch('http://localhost:8000/api/updateProfileImage', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`
                },
                body: formData
            });

            if (response.ok) {
                const data = await response.json();
                setUserImageName(data.profileImageName)
                setUserImageUrl(userImageDirectoryPath + userImageName)
                alert('Image uploaded successfully');
                window.location.reload();
            } else {
                console.error('Failed to upload image');
                alert('Failed to upload image');
            }
        } catch (error) {
            console.error('Error uploading image:', error);
            alert('Error uploading image');
        }
    };

    return(
        <div>
            <NavbarAfterLoggedIn></NavbarAfterLoggedIn>
            <div className="mt-16 p-8 flex">
                {/*User Profile Photo Part */}
                <div className="w-1/2 flex flex-col justify-center items-center border border-gray-500 m-4 rounded-xl">
                <div className="border border-gray-500 w-1/3 h-1/2 rounded-full mx-auto">
                    <img src={userImageUrl} alt={userImageName} className="w-full h-full object-cover rounded-full"></img>
                </div>
                    <form onSubmit={editUserPhoto}>
                        <input type="file" onChange={handleImageChange} className="hover:bg-blue-500 text-white font-bold py-2 px-4 mt-4 rounded cursor-pointer"/>
                        <button type="submit" className="border border-gray-500 p-2 ml-2 rounded hover:bg-blue-500 hover:text-white">Upload</button>
                    </form>
                    {/*<div className="border border-gray-500 w-1/3 h-1/2 rounded-full mx-auto"></div>
                    <button className=" mx-auto p-2 m-2 mt-8 border border-gray-400 rounded" onClick={editUserPhoto}>Edit</button>*/}
                </div>
                {/* User Info Part */}
                <div className="w-1/2 flex flex-col justify-center items-center border border-gray-500 m-4 rounded-xl p-4">
                    <p className="text-center text-xl">User Informations</p>
                    <div className="mt-8">
                        <form onSubmit={handleUpdate} className="flex flex-col">
                            
                            <label htmlFor="firstName" className="block text-sm font-bold mt-2">First Name</label>
                            <input 
                                type="text" id="firstName" name="firstName" value={userData.firstName} placeholder={userData.firstName} onChange={handleChange}
                                className="w-64 mx-auto p-2 border rounded"
                            >
                            </input>

                            <label htmlFor="lastName" className="block text-sm font-bold mt-2">Last Name</label>
                            <input 
                                type="text" id="lastName" name="lastName" value={userData.lastName} placeholder={userData.lastName} onChange={handleChange}
                                className="w-64 mx-auto p-2 border rounded"
                            >
                            </input>

                            <label htmlFor="userName" className="block text-sm font-bold mt-2">User Name</label>
                            <input 
                                type="text" id="userName" name="userName" value={userData.userName} placeholder={userData.userName} onChange={handleChange}
                                className="w-64 mx-auto p-2 border rounded"
                            >
                            </input>

                            <label htmlFor="email" className="block text-sm font-bold mt-2">Email</label>
                            <input 
                                type="text" id="email" name="email" value={userData.email} placeholder={userData.email} onChange={handleChange}
                                className="w-64 mx-auto p-2 border rounded"
                            >
                            </input>
                            
                            <button className=" mx-auto p-2 m-2 mt-8 border border-gray-400 rounded hover:bg-blue-500 hover:text-white">Update</button>

                        </form>
                    </div>
                </div>
            </div>

        </div>
    )
}

export default Profile;