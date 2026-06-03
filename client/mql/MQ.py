import ast
import json

from network.sock import clientCore


class MiQL:
    def __init__(self, __cli : clientCore =None) -> None:
        self.__cli = __cli
        self.__mson = {}
    
    
    @property
    def miql(self):
        return f"mq.{json.dumps(self.__mson)}"
    
    def createDB(self, name, engine="tinydb"):
        self.__mson["active"] = "cre-dat"
        self.__mson["name"] = name
        self.__mson["engine"] = engine
        return self
    
    def dropDB(self, name):
        self.__mson["active"] = "drp-dat"
        self.__mson["name"] = name
        return self

    def shot(self, expr: str | None = None):
        if self.__cli is None:
            raise TypeError("没有设置对应的__cli")

        if expr is not None:
            node = ast.parse(expr, mode="eval").body
            result = self._eval_node(node)
            if isinstance(result, MiQL):
                return result.shot()

            return result

        return self.__cli.send_command(self.miql)

    def _eval_node(self, node):
        
        if isinstance(node, ast.Name):
            if node.id != "mq":
                raise ValueError("MiQL 语句必须以 mq 开头")
            return self

        if isinstance(node, ast.Attribute):
            obj = self._eval_node(node.value)
            if node.attr.startswith("_"):
                raise ValueError(f"不允许访问私有属性: {node.attr}")

            return getattr(obj, node.attr)


        if isinstance(node, ast.Call):
            func = self._eval_node(node.func)
            args = [ast.literal_eval(arg) for arg in node.args]
            kwargs = {
                kw.arg: ast.literal_eval(kw.value)
                for kw in node.keywords
            }
            return func(*args, **kwargs)

        raise ValueError("不支持的 MiQL 语法")

    
    
