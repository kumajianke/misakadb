from typing import List
import random


MAX_LAYER = 32
class PolicyInner:
    def __init__(self, data: SkipListNode | None, poliy_id) -> None:
        self.data = data
        self.id = poliy_id

class Policy:
    GoLowLayer = PolicyInner(None, "GoLowLayer")
    NoFount = PolicyInner(None, "NoFount")
    Empty = PolicyInner(None, "Empty")
    HitIt = PolicyInner(None, "HitIt")
    OnlySmall = PolicyInner(None , "Onlysmall")

class SkipListNode:
    def __init__(self, key, value, next_node=None) -> None:
        self.key = key
        self.value = value
        self.next_node : SkipListNode = next_node
        self.layer_fast_point : List[SkipListNode] = [None for i in range(MAX_LAYER)]


class SkipList:
    def __init__(self):
        self.layer : List[SkipListNode] = [SkipListNode(None, "Header") for i in range(MAX_LAYER)]

    def __random_layer(self):
        layer = 0  # 默认冒泡最底层
        while random.randint(0, 1) == 1 and layer < len(self.layer):
            layer += 1
        return layer

    def insert(self, key, value): 
        random_layer = self.__random_layer # 随机确认这个表冒泡属于哪一层
        new_node = SkipListNode(key, value, None)
        find_max_layer = -1

        for current_layer in range(MAX_LAYER, -1, -1):
            layer : SkipListNode  = self.layer[current_layer]
            policy = self.find_layer_case(layer)
            
            if random_layer > current_layer:
                if policy.id == Policy.Empty.id:
                    layer.next_node = new_node
                if policy.id == Policy.HitIt.id:
                    node : SkipListNode = policy.data.layer_fast_point
                    # 本层级的插入
                    node_next_node = node.next_node
                    node.next_node = new_node
                    new_node.next_node = node_next_node
                    find_max_layer = current_layer
                    break
                
                if current_layer == 0: # 找到最底层
                    if policy.id == Policy.OnlySmall.id:
                        node : SkipListNode = Policy.OnlySmall.data
                        node.next_node = new_node


        for current_layer in range(find_max_layer, -1 , -1):
            node_detail = Policy.GoLowLayer
            

    def find_layer_case(self, layer_linked:SkipListNode, find_node:SkipListNode) -> PolicyInner:
        node = layer_linked.next_node # 头节点
        policy = False

        if node.next_node is None: 
            # 除了头节点什么都没有了
            return Policy.Empty

        while node.next_node: # 头节点的下一个节点 
            next_node = node.next_node
            node_key = node.key

            if node_key is None: # 头节点到下一个节点之间 数据可插入
                policy = find_node.key < next_node.key

            elif node_key == find_node.key or next_node.key == find_node.key:
                Policy.HitIt.data = node_key if node_key == find_node.key else next_node.key
                return Policy.HitIt
            
            elif node_key < find_node.key and next_node is None:
                Policy.OnlySmall.data = {"last": node}
                return Policy.OnlySmall
                
            else: # 节点到节点之间数据可插入
                policy = node_key < find_node.key < next_node.key

            if policy:# 找到数据所在区间 但是没有找到数据
                Policy.GoLowLayer.data = {"last": node , "next": next_node}
                return Policy.GoLowLayer


