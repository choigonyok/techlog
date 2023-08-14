import "./App.css";
import Homepage from "./Homepage/Homepage";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Postpage from "./Postpage/Postpage";
import Writepage from "./Adminpage/Writepage";
import Deletepage from "./Adminpage/Deletepage";
import Loginpage from "./Adminpage/Loginpage";
import axios from "axios";

const App = () => {
  

  return (
    <BrowserRouter>
      <Routes>
        <Route exact path="/" element={<Homepage/>} />
        <Route path="/post/:postid" element={<Postpage/>} />
        <Route path="/admin/delete" element={<Deletepage/>} />
        <Route path="/admin/write" element={<Writepage/>} />
        <Route path="/admin" element={<Loginpage/>} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;

//css는 아래쪽에 있는 코드를 우선순위로 생각한다.
//클래스네임 여러개는  "class1 class2" 로 표현
