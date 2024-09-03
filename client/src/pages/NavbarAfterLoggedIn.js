import {Link} from 'react-router-dom'

const NavbarAfterLoggedIn = () => {
    const handleLogOut = () => {
        console.log("logged out")
        window.location.href = '/';
        localStorage.removeItem("jwtToken");
        alert('Logout successful!');
    }
    return(
        <div>
            <ul className='flex justify-end'>
                <li className='p-4'>
                    <Link className='hover:bg-green-600 p-2 rounded-lg' to="/">Home</Link>
                </li>
                <li className='p-4'>
                    <Link className='hover:bg-green-600 p-2 rounded-lg' to="/chat">Chat</Link>
                </li>
                <li className='p-4'>
                    <Link className='hover:bg-green-600 p-2 rounded-lg' to="/findAFriend">Find a Friend</Link>
                </li>
                <li className='p-4'>
                    <Link className='hover:bg-green-600 p-2 rounded-lg' to="/profile">Profile</Link>
                </li>
                <li className='p-4 cursor-pointer'>
                    <Link onClick={handleLogOut} className='hover:bg-green-600 p-2 rounded-lg'> Log Out</Link>
                </li >
            </ul>
        </div>
    )
}

export default NavbarAfterLoggedIn;