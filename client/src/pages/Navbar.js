import {Link} from 'react-router-dom'

const Navbar = () => {
    return(
        <div>
            <ul className='flex justify-end'>
                <li className='p-4'>
                    <Link className='hover:bg-green-600 p-2 rounded-lg' to="/">Home</Link>
                </li>
                <li className='p-4'>
                    <Link className='hover:bg-green-600 p-2 rounded-lg' to="/login">Login</Link>
                </li>
                <li className='p-4'>
                    <Link className='hover:bg-green-600 p-2 rounded-lg' to="/signup">Sign Up</Link>
                </li >
            </ul>
        </div>
    )
}

export default Navbar;