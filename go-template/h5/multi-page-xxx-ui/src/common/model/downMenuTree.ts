export interface DownMenuTree {
    id: number
    name: string
    parentId: number
    child: DownMenuTree[]|null
    divider?: boolean
    icon: string,
    code: string,
}
