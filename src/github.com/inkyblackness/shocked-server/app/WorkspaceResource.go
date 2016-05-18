package app

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"image/color"
	"image/png"

	"github.com/emicklei/go-restful"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/image"
	core "github.com/inkyblackness/shocked-core"
	model "github.com/inkyblackness/shocked-model"
)

// WorkspaceResource handles all requests for a workspace.
type WorkspaceResource struct {
	ws *core.Workspace
}

// NewWorkspaceResource returns a new workspace resource instance.
func NewWorkspaceResource(container *restful.Container, workspace *core.Workspace) *WorkspaceResource {
	resource := &WorkspaceResource{
		ws: workspace}

	service1 := new(restful.WebService)

	service1.
		Path("/ws").
		Doc("Manage workspace").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service1.Route(service1.GET("").To(resource.getWorkspace).
		// docs
		Doc("get current workspace").
		Operation("getWorkspace").
		Writes(model.Workspace{}))

	container.Add(service1)

	service2 := new(restful.WebService)

	service2.
		Path("/projects").
		Doc("Manage projects").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service2.Route(service2.GET("").To(resource.getProjects).
		// docs
		Doc("get current projects").
		Operation("getWorkspace").
		Writes(model.Projects{}))

	service2.Route(service2.POST("").To(resource.createProject).
		// docs
		Doc("create a project").
		Operation("createProject").
		Reads(model.ProjectTemplate{}).
		Writes(model.Project{}))

	service2.Route(service2.GET("{project-id}/palettes/{palette-id}").To(resource.getPalette).
		// docs
		Doc("get palette").
		Operation("getPalette").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("palette-id", "identifier of the palette").DataType("string")).
		Writes(model.Palette{}))

	service2.Route(service2.GET("{project-id}/textures").To(resource.getTextures).
		// docs
		Doc("get textures").
		Operation("getTextures").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Writes(model.Textures{}))

	service2.Route(service2.GET("{project-id}/textures/{texture-id}").To(resource.getTexture).
		// docs
		Doc("get texture").
		Operation("getTexture").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Writes(model.Texture{}))

	service2.Route(service2.PUT("{project-id}/textures/{texture-id}").To(resource.setTexture).
		// docs
		Doc("set texture").
		Operation("setTexture").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Reads(model.TextureProperties{}).
		Writes(model.Texture{}))

	service2.Route(service2.GET("{project-id}/textures/{texture-id}/{texture-size}").To(resource.getTextureImage).
		// docs
		Doc("get texture image").
		Operation("getTextureImage").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Param(service2.PathParameter("texture-size", "Size of the texture").DataType("string")).
		Writes(model.Image{}))

	service2.Route(service2.GET("{project-id}/textures/{texture-id}/{texture-size}/raw").To(resource.getTextureImageAsRaw).
		// docs
		Doc("get texture image as raw bitmap").
		Operation("getTextureImageAsRaw").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Param(service2.PathParameter("texture-size", "Size of the texture").DataType("string")).
		Writes(model.RawBitmap{}))

	service2.Route(service2.GET("{project-id}/textures/{texture-id}/{texture-size}/png").To(resource.getTextureImageAsPng).
		// docs
		Doc("get texture image as PNG").
		Operation("getTextureImageAsPng").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Param(service2.PathParameter("texture-size", "Size of the texture").DataType("string")).
		Produces("image/png"))

	service2.Route(service2.GET("{project-id}/objects/{class}/{subclass}/{type}").To(resource.getGameObject).
		// docs
		Doc("get game object").
		Operation("getGameObject").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("class", "identifier of the class").DataType("int")).
		Param(service2.PathParameter("subclass", "identifier of the class").DataType("int")).
		Param(service2.PathParameter("type", "identifier of the class").DataType("int")).
		Writes(model.GameObject{}))

	service2.Route(service2.GET("{project-id}/archive/levels").To(resource.getLevels).
		// docs
		Doc("get level list").
		Operation("getLevels").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Writes(model.Levels{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}").To(resource.getLevel).
		// docs
		Doc("get level information").
		Operation("getLevel").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Writes(model.Level{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}/textures").To(resource.getLevelTextures).
		// docs
		Doc("get level textures").
		Operation("getLevelTextures").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Writes(model.LevelTextures{}))

	service2.Route(service2.PUT("{project-id}/archive/levels/{level-id}/textures").To(resource.setLevelTextures).
		// docs
		Doc("put level textures").
		Operation("setLevelTextures").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Reads([]int{}).
		Writes(model.LevelTextures{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}/tiles").To(resource.getLevelTiles).
		// docs
		Doc("get level tiles").
		Operation("getLevelTiles").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Writes(model.Tiles{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}/tiles/{y}/{x}").To(resource.getLevelTile).
		// docs
		Doc("get level tile").
		Operation("getLevelTile").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Param(service2.PathParameter("y", "Y coordinate of the tile").DataType("int")).
		Param(service2.PathParameter("x", "X coordinate of the tile").DataType("int")).
		Writes(model.Tile{}))

	service2.Route(service2.PUT("{project-id}/archive/levels/{level-id}/tiles/{y}/{x}").To(resource.setLevelTile).
		// docs
		Doc("set level tile").
		Operation("setLevelTile").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Param(service2.PathParameter("y", "Y coordinate of the tile").DataType("int")).
		Param(service2.PathParameter("x", "X coordinate of the tile").DataType("int")).
		Reads(model.TileProperties{}).
		Writes(model.Tile{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}/objects").To(resource.getLevelObjects).
		// docs
		Doc("get level objects").
		Operation("getLevelObjects").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Writes(model.LevelObjects{}))

	container.Add(service2)

	return resource
}

// GET /ws
func (resource *WorkspaceResource) getWorkspace(request *restful.Request, response *restful.Response) {
	var entity model.Workspace
	entity.Href = "/"

	entity.Projects.Href = "/projects"

	response.WriteEntity(entity)
}

// GET /projects
func (resource *WorkspaceResource) getProjects(request *restful.Request, response *restful.Response) {
	projectNames := resource.ws.ProjectNames()
	var entity model.Projects
	entity.Href = "/projects"

	entity.Items = make([]model.Identifiable, len(projectNames))
	for index, name := range projectNames {
		proj := &entity.Items[index]
		proj.ID = name
		proj.Href = entity.Href + "/" + proj.ID
	}

	response.WriteEntity(entity)
}

// POST /projects
func (resource *WorkspaceResource) createProject(request *restful.Request, response *restful.Response) {
	entityTemplate := new(model.ProjectTemplate)
	err := request.ReadEntity(entityTemplate)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, prjErr := resource.ws.NewProject(entityTemplate.ID)
	if prjErr != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	entity := new(model.Project)
	entity.ID = entityTemplate.ID
	entity.Href = "/projects/" + entity.ID

	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(entity)
}

// GET /projects/{project-id}/palettes/{palette-id}
func (resource *WorkspaceResource) getPalette(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		paletteID := request.PathParameter("palette-id")
		var palette color.Palette

		palette, err = project.Palettes().GamePalette()

		if paletteID == "game" && err == nil {
			var entity model.Palette

			resource.encodePalette(&entity.Colors, palette)
			response.WriteEntity(entity)
		} else {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusBadRequest, "Unknown palette")
		}
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) encodePalette(out *[256]model.Color, palette color.Palette) {
	for index, inColor := range palette {
		outColor := &out[index]
		r, g, b, _ := inColor.RGBA()

		outColor.Red = int(r >> 8)
		outColor.Green = int(g >> 8)
		outColor.Blue = int(b >> 8)
	}
}

// GET /projects/{project-id}/textures
func (resource *WorkspaceResource) getTextures(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		textures := project.Textures()
		limit := textures.TextureCount()
		var entity model.Textures

		entity.List = make([]model.Texture, limit)
		for id := 0; id < limit; id++ {
			entity.List[id] = resource.textureEntity(project, id)
		}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/textures/{texture-id}
func (resource *WorkspaceResource) getTexture(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		textureID, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		entity := resource.textureEntity(project, int(textureID))

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// PUT /projects/{project-id}/textures/{texture-id}
func (resource *WorkspaceResource) setTexture(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		textureID, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		var properties model.TextureProperties
		err = request.ReadEntity(&properties)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		project.Textures().SetProperties(int(textureID), properties)
		entity := resource.textureEntity(project, int(textureID))

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) textureEntity(project *core.Project, textureID int) (entity model.Texture) {
	entity.ID = fmt.Sprintf("%d", textureID)
	entity.Href = "/projects/" + project.Name() + "/textures/" + entity.ID
	entity.Properties = project.Textures().Properties(textureID)
	for _, size := range model.TextureSizes() {
		entity.Images = append(entity.Images, model.Link{Rel: string(size), Href: entity.Href + "/" + string(size)})
	}

	return
}

// GET /projects/{project-id}/textures/{texture-id}/{texture-size}
func (resource *WorkspaceResource) getTextureImage(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		textureID, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		textureSize := request.PathParameter("texture-size")
		var entity model.Image

		entity.Href = "/projects/" + projectID + "/textures/" + fmt.Sprintf("%d", textureID) + "/" + textureSize
		bmp := project.Textures().Image(int(textureID), model.TextureSize(textureSize))
		hotspot := bmp.Hotspot()

		entity.Properties.HotspotLeft = hotspot.Min.X
		entity.Properties.HotspotTop = hotspot.Min.Y
		entity.Properties.HotspotRight = hotspot.Max.X
		entity.Properties.HotspotBottom = hotspot.Max.Y

		entity.Formats = []model.Link{model.Link{Rel: "png", Href: entity.Href + "/png"}}
		entity.Formats = []model.Link{model.Link{Rel: "raw", Href: entity.Href + "/raw"}}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/textures/{texture-id}/{texture-size}/raw
func (resource *WorkspaceResource) getTextureImageAsRaw(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		textureID, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		textureSize := request.PathParameter("texture-size")
		bmp := project.Textures().Image(int(textureID), model.TextureSize(textureSize))
		var entity model.RawBitmap

		entity.Width = int(bmp.ImageWidth())
		entity.Height = int(bmp.ImageHeight())
		var pixel []byte

		for row := 0; row < entity.Height; row++ {
			pixel = append(pixel, bmp.Row(row)...)
		}
		entity.Pixel = base64.StdEncoding.EncodeToString(pixel)

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/textures/{texture-id}/{texture-size}/png
func (resource *WorkspaceResource) getTextureImageAsPng(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		textureID, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		textureSize := request.PathParameter("texture-size")
		var palette color.Palette

		bmp := project.Textures().Image(int(textureID), model.TextureSize(textureSize))
		palette, err = project.Palettes().GamePalette()
		image := image.FromBitmap(bmp, palette)

		response.AddHeader("Content-Type", "image/png")
		png.Encode(response.ResponseWriter, image)

	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/archive/levels
func (resource *WorkspaceResource) getLevels(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		var entity model.Levels
		archive := project.Archive()
		levelIDs := archive.LevelIDs()

		entity.Href = "/projects/" + projectID + "/archive/levels"
		for _, id := range levelIDs {
			entry := resource.getLevelEntity(project, archive, id)

			entity.List = append(entity.List, entry)
		}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/archive/levels/{level-id}
func (resource *WorkspaceResource) getLevel(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		levelID, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		entity := resource.getLevelEntity(project, project.Archive(), int(levelID))

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) getLevelEntity(project *core.Project, archive *core.Archive, levelID int) (entity model.Level) {
	entity.ID = fmt.Sprintf("%d", levelID)
	entity.Href = "/projects/" + project.Name() + "/archive/levels/" + entity.ID
	level := archive.Level(levelID)
	entity.Properties = level.Properties()

	entity.Links = []model.Link{}
	entity.Links = append(entity.Links, model.Link{Rel: "tiles", Href: entity.Href + "/tiles/{y}/{x}"})
	if !entity.Properties.CyberspaceFlag {
		entity.Links = append(entity.Links, model.Link{Rel: "textures", Href: entity.Href + "/textures"})
	}

	return
}

// GET /projects/{project-id}/archive/levels/{level-id}/textures
func (resource *WorkspaceResource) getLevelTextures(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		levelID, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelID))
		entity := resource.getLevelTexturesEntity(projectID, level)

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) getLevelTexturesEntity(projectID string, level *core.Level) (entity model.LevelTextures) {
	entity.Href = "/projects/" + projectID + "/archive/levels/" + fmt.Sprintf("%d", level.ID()) + "/textures"
	entity.IDs = level.Textures()

	return
}

// PUT /projects/{project-id}/archive/levels/{level-id}/textures
func (resource *WorkspaceResource) setLevelTextures(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		levelID, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)

		var ids []int
		err = request.ReadEntity(&ids)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		level := project.Archive().Level(int(levelID))
		level.SetTextures(ids)

		entity := resource.getLevelTexturesEntity(projectID, level)
		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/archive/levels/{level-id}/tiles
func (resource *WorkspaceResource) getLevelTiles(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		levelID, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelID))
		var entity model.Tiles

		entity.Table = make([][]model.Tile, 64)
		for y := 0; y < 64; y++ {
			entity.Table[y] = make([]model.Tile, 64)
			for x := 0; x < 64; x++ {
				entity.Table[y][x] = getLevelTileEntity(project, level, x, y)
			}
		}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func getLevelTileEntity(project *core.Project, level *core.Level, x int, y int) (entity model.Tile) {
	entity.Href = "/projects/" + project.Name() + "/archive/levels/" + fmt.Sprintf("%d", level.ID()) +
		"/tiles/" + fmt.Sprintf("%d", y) + "/" + fmt.Sprintf("%d", x)
	entity.Properties = level.TileProperties(int(x), int(y))

	return
}

// GET /projects/{project-id}/archive/levels/{level-id}/tiles/{y}/{x}
func (resource *WorkspaceResource) getLevelTile(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		x, _ := strconv.ParseInt(request.PathParameter("x"), 10, 16)
		y, _ := strconv.ParseInt(request.PathParameter("y"), 10, 16)
		levelID, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelID))

		response.WriteEntity(getLevelTileEntity(project, level, int(x), int(y)))
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// PUT /projects/{project-id}/archive/levels/{level-id}/tiles/{y}/{x}
func (resource *WorkspaceResource) setLevelTile(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		x, _ := strconv.ParseInt(request.PathParameter("x"), 10, 16)
		y, _ := strconv.ParseInt(request.PathParameter("y"), 10, 16)
		levelID, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelID))

		var properties model.TileProperties
		err = request.ReadEntity(&properties)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		level.SetTileProperties(int(x), int(y), properties)
		response.WriteEntity(getLevelTileEntity(project, level, int(x), int(y)))
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/archive/levels/{level-id}/objects
func (resource *WorkspaceResource) getLevelObjects(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		levelID, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelID))
		hrefBase := "/projects/" + projectID + "/archive/levels/" + fmt.Sprintf("%d", levelID) + "/objects/"
		var entity model.LevelObjects

		entity.Table = level.Objects()
		for i := 0; i < len(entity.Table); i++ {
			entry := &entity.Table[i]
			entry.Href = hrefBase + entry.ID

			entry.Links = append(entry.Links, model.Link{
				Rel:  "static",
				Href: "/projects/" + projectID + "/objects/" + fmt.Sprintf("%d/%d/%d", entry.Class, entry.Subclass, entry.Type)})
		}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/objects/{class}/{subclass}/{type}
func (resource *WorkspaceResource) getGameObject(request *restful.Request, response *restful.Response) {
	projectID := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectID)

	if err == nil {
		classID, _ := strconv.ParseInt(request.PathParameter("class"), 10, 8)
		subclassID, _ := strconv.ParseInt(request.PathParameter("subclass"), 10, 8)
		typeID, _ := strconv.ParseInt(request.PathParameter("type"), 10, 8)
		objID := res.MakeObjectID(res.ObjectClass(classID), res.ObjectSubclass(subclassID), res.ObjectType(typeID))
		entity := resource.objectEntity(project, objID)

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) objectEntity(project *core.Project, objID res.ObjectID) (entity model.GameObject) {
	entity.ID = fmt.Sprintf("%d/%d/%d", objID.Class, objID.Subclass, objID.Type)
	entity.Href = "/projects/" + project.Name() + "/objects/" + entity.ID
	entity.Properties = project.GameObjects().Properties(objID)

	return
}
