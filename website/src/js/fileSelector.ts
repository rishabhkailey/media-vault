export const selectFile : (callback: ()=>{}) =  (callback) => {
  const inputElement = document.createElement("input");
  inputElement.type = "file";
  inputElement.click();

  return []
}