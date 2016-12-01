package ogl

func NewBox(width, height, depth, textureRepeat float32, shader *Shader) *Model {
	ret := new(Model)

	// Create vertecies etc
	if height == 0 {
		// Just a simple plane no sides
		ret.vertecies = append(ret.vertecies, -width, 0, -depth)
		ret.vertecies = append(ret.vertecies, width, 0, -depth)
		ret.vertecies = append(ret.vertecies, width, 0, depth)
		ret.vertecies = append(ret.vertecies, -width, 0, depth)

		ret.texCoords = append(ret.texCoords, 0, 0)
		ret.texCoords = append(ret.texCoords, textureRepeat, 0)
		ret.texCoords = append(ret.texCoords, textureRepeat, textureRepeat)
		ret.texCoords = append(ret.texCoords, 0, textureRepeat)

		ret.normals = append(ret.normals, 0, 1, 0)
		ret.normals = append(ret.normals, 0, 1, 0)
		ret.normals = append(ret.normals, 0, 1, 0)
		ret.normals = append(ret.normals, 0, 1, 0)

		ret.index = append(ret.index, 0, 1, 2, 0, 2, 3)

	}

	ret.build(shader)

	return ret
}
