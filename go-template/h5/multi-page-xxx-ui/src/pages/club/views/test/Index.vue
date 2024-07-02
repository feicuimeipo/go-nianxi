
<template>
	<div id="editor-builder">
		<div id="top-panel">
		
		
		</div>

		<div id="left-panel">
		</div>

		<div id="canvas">
		</div>

		<div id="right-panel">
		</div>

		<div id="bottom-panel">
		</div>

	</div>

	<!--
		

		

		<div id="top-panel"> 
			<LogoLeft />
			<CenterBtn />
			<RightBtn />        
		</div>

		<LeftPanel />
		<CanvasPanel />

		<div id="right-panel">
			<div class="component-properties">
				<ul class="nav nav-tabs nav-fill" id="properties-tabs" role="tablist">
					<RightPanelTabs />
				</ul>

				<div class="tab-component_block">
					<RightPanelTabContent/>
				</div>           
			
			</div>
		</div>
		
		<BottomPanel />
		<FormTemplate /> 
		<FullTemplate />
		<HtmlModal />
		
		<SaveToast />
	</div> -->
	

</template>

<script lang="ts">
	// <!--top-->
	// import LogoBtn from  '@/components/editor/top/LogoLeft.vue';
	// import CenterBtn from '@/components/editor/top/BtnCenter.vue';	
	// import RightBtn from '@/components/editor/top/BtnRight.vue';
	// import LeftPanel from "@/components/editor/LeftPanel.vue"
	// import CanvasPanel from "@/components/editor/CanvasPanel.vue"
	// import BottomPanel from "@/components/editor/BottomPanel.vue"
	// import HtmlModal from "@/components/editor/HtmlModal.vue"
	// import FullTemplate from "@/components/editor/template/FullTemplate.vue"	
	// import SaveToast from "@/components/editor/SaveToast.vue"
	// import RightPanelTabContent from '@/components/editor/right/RightPanelTabContent.vue'
	// import RightPanelTabs from '@/components/editor/right/RightPanelTabs.vue'
	
	import { nextTick } from "vue"		
	import $ from 'jquery' 	

	
 
	export default {
		name: 'Test',
		mounted:function () {		
			$(function() {			
				//alert("页面已加载")	
			});
		}
  } 

</script>


<style  scoped>
</style>
