package kubernetes
import("gopkg.in/yaml.v3";"strings")
func Pin(b []byte,res func(string)(string,error))([]byte,bool,error){
 var obj map[string]interface{}
 if err:=yaml.Unmarshal(b,&obj); err!=nil{return b,false,nil}
 changed:=false; walk(obj,res,&changed)
 out,_:=yaml.Marshal(obj); return out,changed,nil
}
func walk(n interface{},res func(string)(string,error),c *bool){
 switch v:=n.(type){
 case map[string]interface{}:
   for k,val:=range v{
    if k=="image"{
      if img,ok:=val.(string);ok&&!strings.Contains(img,"@sha256:"){
        if dg,err:=res(img);err==nil{v[k]=dg;*c=true}
      }
    }
    walk(val,res,c)
   }
 case []interface{}: for _,it:=range v{walk(it,res,c)}
 }
}
