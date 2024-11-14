import "./App.css";
import Homepage from "./Homepage/Homepage";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Postpage from "./Postpage/Postpage";
import Writepage from "./Adminpage/Writepage";
import Deletepage from "./Adminpage/Deletepage";
import Login from "./Header/Login";
import Callback from "./Header/Callback";

const App = () => {
  

  return (
    <BrowserRouter>
      <Routes>
        <Route exact path="/" element={<Homepage/>} />
        <Route path="/posts/:postid" element={<Postpage/>} />
        <Route path="/admin" element={<Deletepage/>} />
        <Route path="/login" element={<Login/>} />
        <Route path="/github/callback" element={<Callback/>} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;