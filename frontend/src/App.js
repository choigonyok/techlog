import "./App.css";
import Homepage from "./Homepage/Homepage";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Postpage from "./Postpage/Postpage";
import Writepage from "./Adminpage/Writepage";
import Deletepage from "./Adminpage/Deletepage";
import Loginpage from "./Adminpage/Loginpage";

const App = () => {
  

  return (
    <BrowserRouter>
      <Routes>
        <Route exact path="/" element={<Homepage/>} />
        <Route path="/post/:postid" element={<Postpage/>} />
        <Route path="/admin/delete" element={<Deletepage/>} />
        <Route path="/admin/write" element={<Writepage/>} />
        <Route path="/login" element={<Loginpage/>} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;