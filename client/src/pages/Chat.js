import NavbarAfterLoggedIn from "./NavbarAfterLoggedIn"
import { useState, useEffect, useRef } from "react";
import backgroundImage from '../images/wpBg.png';

const Chat = () =>{

    //

    const userImageDirectoryPath = "http://127.0.0.1:8000/userImages/";

    const [friends, setFriends] = useState([]);

    const [ws, setWs] = useState(null);
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState([]);
    const [groupName, setGroupName] = useState("");
    const [userName, setUserName] = useState("");
    const [userId, setUserId] = useState(0);

    const [selectedFriendId, setSelectedFriendId] = useState(0)
    const [selectedFriendFirstName, setSelectedFriendFirstName] = useState("")
    const [selectedFriendLastName, setSelectedFriendLastName] = useState("")
    const [selectedFriendUserName, setSelectedFriendUserName] = useState("")
    const [selectedFriendImageName, setSelectedFriendImageName] = useState("")

    const messagesEndRef = useRef(null);
    const friendsEndRef = useRef(null);

    const [checkNewMessagesArray, setCheckNewMessagesArray] = useState({})
    const [newMessages, setNewMessages] = useState({})

    useEffect(() => {
        // Scroll to the bottom whenever messages change
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
        }
        if (friendsEndRef.current) {
            friendsEndRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [messages,friends]);

    const fetchUserFriendsData = async () => {
        try{
            const token = localStorage.getItem("jwtToken");

            const response = await fetch('http://localhost:8000/api/getFriends', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                // Token is valid, extract user data from response
                const responseData = await response.json();
                setFriends(responseData.friends);
                
            } else {
                console.error('Couldnt fetch user data');
            }

        }
        catch (error) {
            console.error('Error fetching user data :', error);
        }
    }

    const fetchUserData = async () => {
        try{
            const token = localStorage.getItem("jwtToken");

            const response = await fetch('http://localhost:8000/api/getUserName', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                // Token is valid, extract user data from response
                const responseData = await response.json();
                setUserName(responseData.userName);
                setUserId(responseData.id);
                
            } else {
                console.error('Couldnt fetch user name');
            }

        }
        catch (error) {
            console.error('Error fetching user name :', error);
        }
    }

    const checkNewMessages = async () => {
        try{
            const token = localStorage.getItem("jwtToken");

            const response = await fetch('http://localhost:8000/api/CheckMessages', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                // Token is valid, extract user data from response
                const responseData = await response.json();
                setCheckNewMessagesArray(responseData.data)

                
            } else {
                console.error('Couldnt check new messages ');
            }

        }
        catch (error) {
            console.error('Error checking new messages :', error);
        }
    }

    const getNewMessages = async (groupName, userName, selectedFriendId, userId) => {
        try{
            const token = localStorage.getItem("jwtToken");

            const response = await fetch(`http://localhost:8000/api/getNewMessages?Group=${groupName}&User=${userName}&SelectedFriendId=${selectedFriendId}&UserId=${userId}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                // Token is valid, extract user data from response
                const responseData = await response.json();
                setNewMessages(responseData.data)

                
            } else {
                console.error('Couldnt check new messages ');
            }

        }
        catch (error) {
            console.error('Error checking new messages :', error);
        }
    }

    const generateGroupString = () => {
        if (selectedFriendId && userId) {
            // Ensure both IDs are strings
            const id1 = String(selectedFriendId);
            const id2 = String(userId);

            // Generate group string
            const groupName = [id1, id2].sort().join("");
            setGroupName(groupName);
        }
    };

    //check new messages
    useEffect(()=>{
        checkNewMessages()
    },[selectedFriendId]);
    
    //get friends
    useEffect(()=>{

        fetchUserFriendsData();
        fetchUserData();



        if( selectedFriendId !== 0 ){

            setMessages([])

            generateGroupString();

            const ws = new WebSocket(`ws://localhost:8000/ws?Group=${groupName}&User=${userName}&SelectedFriendId=${selectedFriendId}&UserId=${userId}`);

            ws.onopen = function() {
                console.log("Connected to WebSocket server");
                getNewMessages(groupName, userName, selectedFriendId, userId);
            };
            ws.onmessage = function(event) {
                setMessages(prevMessages => [...prevMessages, [1,event.data]]);
                console.log("Message from server:", event.data);
            };
            ws.onclose = function() {
                console.log("Disconnected from WebSocket server");
            };
            ws.onerror = function(error) {
                console.log("WebSocket error:", error);
            };
    
            setWs(ws);
    
            return () => ws.close();            
        }

    }, [selectedFriendId,userId,groupName,userName]);

    const selectFriend = (uid, fn, ls, un, pin) =>{
        setSelectedFriendId(uid)
        setSelectedFriendFirstName(fn)
        setSelectedFriendLastName(ls)
        setSelectedFriendUserName(un)
        setSelectedFriendImageName(pin)
    }

    const handleSendMessage = () => {
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(message);
            // 0 mean the message is sent by user, 1 means message is sent by friend
            setMessages(prevMessages => [...prevMessages, [0,message]]);
            setMessage(""); // Clear the textarea after sending
        } else {
            console.log("WebSocket is not open.");
        }
    };

    return(
        <div>
            <NavbarAfterLoggedIn></NavbarAfterLoggedIn>
            <div className="flex">
                {/*Users In Contact*/}
                <div className="w-1/4 h-[85vh] border m-4 border-gray-500 rounded-lg" style={{backgroundColor: '#334155'}}>
                    <p className="p-2 m-2 rounded-lg border-b border-gray-500 bg-green-600 text-center text-white">Messages</p>
                    {friends &&
                        friends.map(friend => (
                            <div key={friend.id} style={{backgroundColor: checkNewMessagesArray[friend.id] ? '#ecfdf5' : '#334155' }} className="flex items-center p-2 m-2 rounded-md border-b border-gray-500 cursor-pointer" onClick={()=>selectFriend(friend.id, friend.firstName, friend.lastName, friend.userName, friend.profileImageName)}>
                                <div className="border border-gray-500 w-12 h-12 rounded-full mr-2">
                                    <img src={userImageDirectoryPath + friend.profileImageName} alt={friend.profileImageName} className="w-full h-full object-cover rounded-full"></img>
                                </div>
                                <div>
                                    <p className="text-sm font-bold">{friend.firstName} {friend.lastName}</p>
                                    <p className="text-sm">{friend.userName}</p>
                                </div>
                                {checkNewMessagesArray[friend.id] && (
                                    <div className="bg-green-500 rounded-lg ml-auto p-1">
                                        <p className="text-sm font-bold">{checkNewMessagesArray[friend.id]}</p>
                                    </div>
                                )}
                            </div>
                        ))
                    }
                    <div ref={friendsEndRef} />
                </div>
                {/*Chat Messages*/}
                {selectedFriendFirstName &&
                    <div className="w-3/4 flex flex-col justify-between border m-1 border-gray-500 rounded-lg">
                        <div className="flex p-2 border-b border-gray-500 bg-green-50" style={{backgroundColor: '#047857'}}>
                            <div className="border border-gray-500 w-12 h-12 rounded-full mr-2">
                                <img src={userImageDirectoryPath + selectedFriendImageName} alt={selectedFriendImageName} className="w-full h-full object-cover rounded-full"></img>
                            </div>
                            <div>
                                <p className="text-sm font-bold text-white">{selectedFriendFirstName} {selectedFriendLastName}</p>
                                <p className="text-sm text-white">{selectedFriendUserName}</p>
                            </div>
                        </div>
                        <div style={{ backgroundImage: `url(${backgroundImage})` }} className=" rounded-lg m-1 h-[70vh] flex flex-col overflow-auto bg-cover bg-center">
                            {messages.map((msg, index) => (
                                <div 
                                    key={index} 
                                    style={{backgroundColor: msg[0] === 0 ? '#047857' : '#334155'}}
                                    className={`message m-2 p-2 max-w-xl h-fit rounded-lg break-words ${msg[0] === 0 ? 'text-white self-end text-left' : 'text-white self-start text-left'}`}
                                >
                                    {msg.slice(1)}
                                </div>
                            ))}
                            <div ref={messagesEndRef} />
                        </div>
                        <div className="flex justify-center items-center p-2" style={{backgroundColor: '#334155'}}>
                            <textarea className="w-4/5 h-12 p-2 mx-2 border rounded-lg focus:outline-none focus:ring focus:border-blue-500 resize-none"value={message} onChange={(e) => setMessage(e.target.value)}></textarea>
                            <button className="p-2 border border-gray-400 rounded text-white hover:bg-blue-500 hover:text-white" style={{backgroundColor: '#047857'}} onClick={handleSendMessage}>Send</button>
                        </div>
                    </div>
                }

            </div>
        </div>
    )
}

export default Chat