import React, { useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import "./Markdown.css"; // CSS 파일 작성

const Markdown = ({onTextChange, value}) => {
  const [markdown, setMarkdown] = useState(value);

  // 입력값 변경 핸들러
  const handleInputChange = (event) => {
    setMarkdown(event.target.value);
    onTextChange(event.target.value)
  };

  return (
    <div className="markdown-editor">
      <div className="editor-container">
        {/* 입력 영역 */}
        <textarea
          className="markdown-input"
          value={markdown}
          onChange={handleInputChange}
          placeholder="Type your markdown here..."
          style={{backgroundColor: "rgb(59, 59, 59)", color: "white"}}
        ></textarea>

        {/* 미리보기 영역 */}
        <div className="markdown-preview">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>
            {markdown}
          </ReactMarkdown>
        </div>
      </div>
    </div>
  );
};

export default Markdown;
