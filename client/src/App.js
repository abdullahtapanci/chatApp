import './App.css';
import {BrowserRouter as Router,Routes, Route} from 'react-router-dom'
import Home from "./pages/Home.js"
import Login from "./pages/Login.js"
import Signup from "./pages/Signup.js"
import Profile from "./pages/Profile.js"
import Chat from './pages/Chat.js';
import FindAFriend from './pages/FindAFriend.js';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home></Home>}></Route>
        <Route path="/login" element={<Login></Login>}></Route>
        <Route path="/signup" element={<Signup></Signup>}></Route>
        <Route path="/profile" element={<Profile></Profile>}></Route>
        <Route path="/chat" element={<Chat></Chat>}></Route>
        <Route path="/findAFriend" element={<FindAFriend></FindAFriend>}></Route>
      </Routes>
    </Router>
  );
}

export default App;
