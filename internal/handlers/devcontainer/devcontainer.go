package devcontainer
import("encoding/json";"strings")
func Pin(b []byte,res func(string)(string,error))([]byte,bool,error){
 var dc map[string]interface{}
 if err:=json.Unmarshal(b,&dc); err!=nil{return b,false,nil}
 changed:=false
 if img,ok:=dc["image"].(string); ok && !strings.Contains(img,"@sha256:"){
   if dg,err:=res(img); err==nil{dc["image"]=dg;changed=true}
 }
 out,_:=json.MarshalIndent(dc,"","  ")
 return out,changed,nil
}
